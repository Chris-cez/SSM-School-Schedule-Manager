package scheduler

import (
	"SSM-School-Schedule-Manager/classes"
	"fmt"
)

// Scheduler resolve conflitos de agendamento usando CSP
type Scheduler struct {
	Turmas         []classes.Turma
	Restricoes     classes.Restricoes
	NumHorarios    int
	GrafoConflitos map[int][]int
}

// NewScheduler cria um novo scheduler
func NewScheduler(turmas []classes.Turma, numHorarios int, restricoes classes.Restricoes) *Scheduler {
	return &Scheduler{
		Turmas:         turmas,
		Restricoes:     restricoes,
		NumHorarios:    numHorarios,
		GrafoConflitos: make(map[int][]int),
	}
}

// ConstructGrafo constrói um grafo de conflitos entre turmas
func (s *Scheduler) ConstructGrafo() map[int][]int {
	grafo := make(map[int][]int)
	professorTurmas := make(map[string][]int)

	// Mapear professores para turmas
	for i, turma := range s.Turmas {
		for _, disciplina := range turma.Disciplinas {
			professorTurmas[disciplina.Professor.Nome] = append(professorTurmas[disciplina.Professor.Nome], i)
		}
	}

	// Criar arestas para turmas que compartilham professor
	for _, turmas := range professorTurmas {
		for i := 0; i < len(turmas); i++ {
			for j := i + 1; j < len(turmas); j++ {
				turmaA := turmas[i]
				turmaB := turmas[j]
				grafo[turmaA] = append(grafo[turmaA], turmaB)
				grafo[turmaB] = append(grafo[turmaB], turmaA)
			}
		}
	}

	s.GrafoConflitos = grafo
	return grafo
}

// Resolver usa backtracking para criar horários respeitando restrições
func (s *Scheduler) Resolver() bool {
	s.ConstructGrafo()

	// Inicializar horários vazios para cada turma
	for i := range s.Turmas {
		s.Turmas[i].Horarios = make([]int, 0, s.NumHorarios)
	}

	// Tentar resolver para cada turma
	for i := range s.Turmas {
		s.Turmas[i].Horarios = make([]int, 0)
		if !s.backtrackTurma(i, 0) {
			return false
		}
	}

	return true
}

// backtrackTurma resolve o agendamento de uma turma específica
func (s *Scheduler) backtrackTurma(turmaIdx int, horarioAtual int) bool {
	turma := &s.Turmas[turmaIdx]

	// Verificar se atendemos a frequência de todas as disciplinas
	if horarioAtual == s.NumHorarios {
		return s.validarFrequencia(turmaIdx)
	}

	// Verificar limite de aulas por dia
	aulasNoDia := horarioAtual % s.Restricoes.AulasporDia

	if aulasNoDia == 0 && horarioAtual > 0 {
		dia := (horarioAtual - 1) / s.Restricoes.AulasporDia
		if dia >= s.Restricoes.DiasAtivos {
			return false
		}
	}

	// Tentar cada disciplina da turma
	for _, disciplina := range turma.Disciplinas {
		// Verificar quantas vezes essa disciplina já foi atribuída
		countDisciplina := s.contarDisciplina(turmaIdx, disciplina.ID)
		frequenciaDesejada := s.Restricoes.FrequenciaDisciplina[disciplina.ID]

		if countDisciplina >= frequenciaDesejada {
			continue
		}

		// Verificar restrições
		if s.podeAtribuir(turmaIdx, horarioAtual, disciplina) {
			turma.Horarios = append(turma.Horarios, disciplina.ID)

			if s.backtrackTurma(turmaIdx, horarioAtual+1) {
				return true
			}

			// Backtrack
			turma.Horarios = turma.Horarios[:len(turma.Horarios)-1]
		}
	}

	return false
}

// podeAtribuir verifica se é seguro atribuir uma disciplina a um horário
func (s *Scheduler) podeAtribuir(turmaIdx int, horarioAtual int, disciplina classes.Disciplina) bool {
	turma := &s.Turmas[turmaIdx]

	// Verificar se o professor está disponível neste horário
	if disciplina.Professor.Indisponibilidades != nil && disciplina.Professor.Indisponibilidades[horarioAtual] {
		return false
	}

	// Verificar aulas consecutivas da mesma disciplina
	if len(turma.Horarios) > 0 {
		consecutivas := 1
		for i := len(turma.Horarios) - 1; i >= 0 && turma.Horarios[i] == disciplina.ID; i-- {
			consecutivas++
		}

		if consecutivas > s.Restricoes.MaxAulasConsecutivasDisciplina {
			return false
		}
	}

	// Verificar aulas consecutivas do mesmo professor
	if len(turma.Horarios) > 0 {
		ultimoDisciplinaID := turma.Horarios[len(turma.Horarios)-1]
		ultimaDisciplina := s.encontrarDisciplina(turmaIdx, ultimoDisciplinaID)

		if ultimaDisciplina != nil && ultimaDisciplina.Professor.Nome == disciplina.Professor.Nome {
			consecutivas := 1
			for i := len(turma.Horarios) - 1; i >= 0; i-- {
				discID := turma.Horarios[i]
				disc := s.encontrarDisciplina(turmaIdx, discID)
				if disc != nil && disc.Professor.Nome == disciplina.Professor.Nome {
					consecutivas++
				} else {
					break
				}
			}

			if consecutivas > s.Restricoes.MaxAulasConsecutivasProfessor {
				return false
			}
		}
	}

	// Verificar conflito com outras turmas (professor ocupado)
	if !s.professorDisponivel(horarioAtual, disciplina.Professor.Nome, turmaIdx) {
		return false
	}

	return true
}

// professorDisponivel verifica se o professor está livre neste horário
func (s *Scheduler) professorDisponivel(horarioAtual int, professorNome string, turmaIdx int) bool {
	// Verificar outras turmas que têm aulas neste mesmo horário
	for i, turma := range s.Turmas {
		if i == turmaIdx {
			continue
		}

		if len(turma.Horarios) > horarioAtual {
			disciplinaID := turma.Horarios[horarioAtual]
			for _, disc := range turma.Disciplinas {
				if disc.ID == disciplinaID && disc.Professor.Nome == professorNome {
					return false
				}
			}
		}
	}

	return true
}

// contarDisciplina conta quantas vezes uma disciplina foi atribuída
func (s *Scheduler) contarDisciplina(turmaIdx int, disciplinaID int) int {
	count := 0
	for _, discID := range s.Turmas[turmaIdx].Horarios {
		if discID == disciplinaID {
			count++
		}
	}
	return count
}

// encontrarDisciplina encontra uma disciplina pelo ID
func (s *Scheduler) encontrarDisciplina(turmaIdx int, disciplinaID int) *classes.Disciplina {
	for i, disc := range s.Turmas[turmaIdx].Disciplinas {
		if disc.ID == disciplinaID {
			return &s.Turmas[turmaIdx].Disciplinas[i]
		}
	}
	return nil
}

// validarFrequencia verifica se todas as frequências foram respeitadas
func (s *Scheduler) validarFrequencia(turmaIdx int) bool {
	for disciplinaID, frequenciaDesejada := range s.Restricoes.FrequenciaDisciplina {
		count := s.contarDisciplina(turmaIdx, disciplinaID)
		if count != frequenciaDesejada {
			return false
		}
	}
	return true
}

// Exibir mostra o agendamento
func (s *Scheduler) Exibir() {
	fmt.Println("\n=== AGENDAMENTO FINAL ===\n")

	for _, turma := range s.Turmas {
		fmt.Printf("Turma: %s\n", turma.Nome)

		for i, disciplinaID := range turma.Horarios {
			disc := s.encontrarDisciplina(turma.ID-1, disciplinaID)

			if disc != nil {
				dia := (i / s.Restricoes.AulasporDia) + 1
				aulaDoDia := (i % s.Restricoes.AulasporDia) + 1

				fmt.Printf("  Dia %d - Aula %d: %s - Prof. %s\n",
					dia, aulaDoDia, disc.Nome, disc.Professor.Nome)
			}
		}
		fmt.Println()
	}
}
