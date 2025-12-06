package classes

type Turma struct {
	ID          int
	Nome        string
	Disciplinas []Disciplina
	Horarios    []int
}

type Professor struct {
	ID                 int
	Nome               string
	Indisponibilidades map[int]bool
}

type Disciplina struct {
	ID        int
	Nome      string
	Professor Professor
}

type Restricoes struct {
	MaxAulasConsecutivasDisciplina int
	MaxAulasConsecutivasProfessor  int
	FrequenciaDisciplina           map[int]int
	AulasporDia                    int
	DiasAtivos                     int
}

type Grafo struct {
	Turmas  []Turma
	Arestas map[int][]int
}
