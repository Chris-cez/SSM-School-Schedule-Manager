package classes

type Turma struct {
	ID          int
	Nome        string
	Disciplinas []Disciplina
}

type Professor struct {
	ID    int
	Nome  string
	Email string
}

type Disciplina struct {
	ID        int
	Nome      string
	Professor Professor
}
