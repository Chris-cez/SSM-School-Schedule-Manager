package main

import (
	"SSM-School-Schedule-Manager/classes"
	"SSM-School-Schedule-Manager/scheduler"
	"fmt"
)

func main() {
	// Criar professores com indisponibilidades
	prof1 := classes.Professor{
		ID:   1,
		Nome: "Dr. Silva",
		Indisponibilidades: map[int]bool{
			4: true, // Indisponível no horário 4 (Dia 3 - Aula 1)
		},
	}

	prof2 := classes.Professor{
		ID:   2,
		Nome: "Dra. Santos",
		Indisponibilidades: map[int]bool{
			1: true, // Indisponível no horário 1 (Dia 1 - Aula 2)
			6: true, // Indisponível no horário 6 (Dia 4 - Aula 1)
		},
	}

	prof3 := classes.Professor{
		ID:                 3,
		Nome:               "Prof. Costa",
		Indisponibilidades: make(map[int]bool),
	}

	// Criar disciplinas
	disc1 := classes.Disciplina{ID: 1, Nome: "Matemática", Professor: prof1}
	disc2 := classes.Disciplina{ID: 2, Nome: "Português", Professor: prof2}
	disc3 := classes.Disciplina{ID: 3, Nome: "História", Professor: prof3}

	// Criar turmas
	turma1 := classes.Turma{
		ID:          1,
		Nome:        "Turma A",
		Disciplinas: []classes.Disciplina{disc1, disc2, disc3},
	}
	turma2 := classes.Turma{
		ID:          2,
		Nome:        "Turma B",
		Disciplinas: []classes.Disciplina{disc1, disc2, disc3},
	}

	turmas := []classes.Turma{turma1, turma2}

	// Definir restrições
	restricoes := classes.Restricoes{
		MaxAulasConsecutivasDisciplina: 2, // Máximo 2 aulas seguidas da mesma disciplina
		MaxAulasConsecutivasProfessor:  2, // Máximo 2 aulas seguidas do mesmo professor
		AulasporDia:                    2, // 2 aulas por dia
		DiasAtivos:                     5, // 5 dias úteis na semana
		FrequenciaDisciplina: map[int]int{
			1: 3, // Matemática 3 vezes na semana
			2: 2, // Português 2 vezes na semana
			3: 2, // História 2 vezes na semana
		},
	}

	// Resolver com CSP
	sch := scheduler.NewScheduler(turmas, 10, restricoes) // 10 horários (5 dias * 2 aulas)

	if sch.Resolver() {
		fmt.Println("✓ Agendamento resolvido com sucesso!")
		sch.Exibir()
	} else {
		fmt.Println("✗ Impossível resolver o agendamento com as restrições fornecidas")
	}
}
