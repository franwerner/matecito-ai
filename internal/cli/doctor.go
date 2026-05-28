package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/franwerner/matecito-ai/internal/check"
	"github.com/franwerner/matecito-ai/internal/codegraph"
	"github.com/franwerner/matecito-ai/internal/context7"
	"github.com/franwerner/matecito-ai/internal/engram"
	"github.com/franwerner/matecito-ai/internal/prereqs"
	"github.com/franwerner/matecito-ai/internal/render"
	"github.com/franwerner/matecito-ai/internal/sdd"
)

func NewDoctorCmd() *cobra.Command {
	var sddDir string

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Diagnóstico accionable: sólo problemas + comandos sugeridos",
		Long:  "doctor corre los mismos chequeos que verify, filtra los ✓, agrupa los\ncomandos sugeridos al final para copy-paste, y devuelve el mismo exit code.",
		RunE: func(cmd *cobra.Command, args []string) error {
			sections := []struct {
				title   string
				results []check.Result
			}{
				{"Prerequisites", prereqs.All()},
				{"Engram", engram.All()},
				{"CodeGraph", codegraph.All()},
				{"context7", context7.All()},
				{"Cross-check SDD ↔ MCP (" + sddDir + ")", sdd.CrossCheck(sddDir)},
			}

			var all, problems []check.Result
			anyProblem := false

			for _, sec := range sections {
				all = append(all, sec.results...)
				filtered := filterProblems(sec.results)
				if len(filtered) > 0 {
					render.Section(os.Stdout, sec.title, filtered)
					problems = append(problems, filtered...)
					anyProblem = true
				}
			}

			if !anyProblem {
				fmt.Println("Todo bien — no hay problemas para resolver.")
				return nil
			}

			render.Suggestions(os.Stdout, problems)

			if code := render.Summary(os.Stdout, all); code != 0 {
				os.Exit(code)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&sddDir, "sdd-dir", defaultSDDDir(),
		"Directorio donde viven los agentes del SDD (sdd-*.md)")
	return cmd
}

func filterProblems(results []check.Result) []check.Result {
	out := make([]check.Result, 0, len(results))
	for _, r := range results {
		if r.Status != check.StatusOK {
			out = append(out, r)
		}
	}
	return out
}
