package cli

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func NewInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Inicializa lo por-proyecto en el cwd (.codegraph/)",
		Long:  "init corre `codegraph init -i` en el directorio actual si CodeGraph está\ninstalado y falta `.codegraph/`. No toca config global.",
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("no se pudo resolver cwd: %w", err)
			}

			cgBin, err := exec.LookPath("codegraph")
			if err != nil {
				return errors.New("codegraph no está instalado. Instalalo con: npm install -g @colbymchenry/codegraph")
			}

			cgDir := filepath.Join(cwd, ".codegraph")
			if fi, err := os.Stat(cgDir); err == nil && fi.IsDir() {
				fmt.Printf("✓ .codegraph/ ya existe en %s — nada para hacer.\n", cwd)
				return nil
			}

			fmt.Printf("Inicializando CodeGraph en %s\n", cwd)
			c := exec.Command(cgBin, "init", "-i")
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			c.Dir = cwd
			if err := c.Run(); err != nil {
				return fmt.Errorf("codegraph init falló: %w", err)
			}
			fmt.Println("✓ .codegraph/ inicializado.")
			return nil
		},
	}
}
