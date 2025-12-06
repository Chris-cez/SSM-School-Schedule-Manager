package main

import (
	"classes"
	"fmt"
)

func main() {
	professor := classes.Professor{
		ID:    1,
		Nome:  "Dr. Silva",
		Email: "",
	}
	disciplina := classes.Disciplina{
		ID:        1,
		Nome:      "Matemática",
		Professor: professor,
	}
	turma := classes.Turma{
		ID:          1,
		Nome:        "Turma A",
		Disciplinas: []classes.Disciplina{disciplina},
	}

	fmt.Printf("Turma: %s\n", turma.Nome)
	fmt.Printf("Disciplina: %s\n", turma.Disciplinas[0].Nome)
	fmt.Printf("Professor: %s\n", turma.Disciplinas[0].Professor.Nome)
}
