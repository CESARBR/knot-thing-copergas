package entities

type Affiliation struct {
	CodAgrVin int `json:"codAgrVin"`
	CodAgr    int `json:"codAgr"`
	CodAgrPai int `json:"codAgrPai"`
	CodEst    int `json:"codEst"`
	CodMed    int `json:"codMed"`
	CodVar    int `json:"codVar"`
	CodTgg    int `json:"codTgg"`
	CodEmpr   int `json:"codEmpr"`
}

type Variable struct {
	CodVar                 int           `json:"codVar"`
	CodMed                 int           `json:"codMed"`
	CodEst                 int           `json:"codEst"`
	CodEmpr                int           `json:"codEmpr"`
	Descricao              string        `json:"descricao"`
	Comprimento            int           `json:"comprimento"`
	Pagina                 int           `json:"pagina"`
	Endereco               int           `json:"endereco"`
	PosicaoEnron           string        `json:"posicaoEnron"`
	InicioBloco            bool          `json:"inicioBloco"`
	Unidade                string        `json:"unidade"`
	FatorCorrecao          float32       `json:"fatorCorrecao"`
	ParcelaCorrecao        float32       `json:"parcelaCorrecao"`
	AlmHabilitado          bool          `json:"almHabilitado"`
	CfgAlmMuitoAlto        bool          `json:"cfgAlmMuitoAlto"`
	CfgAlmAlto             bool          `json:"cfgAlmAlto"`
	CfgAlmBaixo            bool          `json:"cfgAlmBaixo"`
	CfgAlmMuitoBaixo       bool          `json:"cfgAlmMuitoBaixo"`
	CfgAlmValor            bool          `json:"cfgAlmValor"`
	CfgAlmMudanca          bool          `json:"CfgAlmMudanca"`
	CfgAlmVariacao         bool          `json:"CfgAlmVariacao"`
	CodTpDado              int           `json:"CodTpDado"`
	CodTpOrdByte           int           `json:"CodTpOrdByte"`
	ValorInteger           int           `json:"ValorInteger"`
	ValorFloat             float32       `json:"ValorFloat"`
	ValorString            string        `json:"ValorString"`
	ValorDateTime          string        `json:"ValorDateTime"`
	ValorConv              interface{}   `json:"ValorConv"`
	ValorConvFormat        string        `json:"ValorConvFormat"`
	ValorConvStrFormat     string        `json:"ValorConvStrFormat"`
	DataLeitura            string        `json:"dataLeitura"`
	EmAlarmeNivel          bool          `json:"emAlarmeNivel"`
	TagWeb                 string        `json:"tagWeb"`
	Tag                    string        `json:"access_token"`
	ValorIntegerEscr       int           `json:"ValorIntegerEscr"`
	ValorFloatEscr         float32       `json:"ValorFloatEscr"`
	ValorStringEscr        string        `json:"ValorStringEscr"`
	ValorDateTimeEscr      string        `json:"ValorDateTimeEscr"`
	ValorConvEscr          string        `json:"ValorConvEscr"`
	ValorConvEscrFormat    string        `json:"ValorConvEscrFormat"`
	ValorConvEscrStrFormat string        `json:"ValorConvEscrStrFormat"`
	CodTpDadoEscr          int           `json:"codTpDadoEscr"`
	CodStExecEscr          int           `json:"codStExecEscr"`
	DataStExecEscr         string        `json:"dataStExecEscr"`
	IntervaloLeituraMin    float32       `json:"intervaloLeituraMin"`
	Prioridade             int           `json:"prioridade"`
	TagWebUn               string        `json:"tagWebUn"`
	CodTgg                 int           `json:"CodTgg"`
	IndicePilhaEnron       int           `json:"IndicePilhaEnron"`
	TamanhoPilhaEnron      int           `json:"TamanhoPilhaEnron"`
	PosicaoEnronUn         int           `json:"PosicaoEnronUn"`
	AhProcesso             string        `json:"AhProcesso"`
	AhSecao                string        `json:"AhSecao"`
	AhItem                 string        `json:"AhItem"`
	FazerTelemetria        bool          `json:"FazerTelemetria"`
	GravarHistorico        bool          `json:"gravarHistorico"`
	EsconderEmDesenho      bool          `json:"esconderEmDesenho"`
	Programar              bool          `json:"programar"`
	CodMedUn               string        `json:"codMedUn"`
	BytesPorRegistro       int           `json:"bytesPorRegistro"`
	FuncaoEscr             string        `json:"funcaoEscr"`
	Escala                 int           `json:"escala"`
	CodFdt                 int           `json:"codFdt"`
	TotalVarBit            string        `json:"totalVarBit"`
	CodEstList             string        `json:"codEstList"`
	BitValue               bool          `json:"bitValue"`
	Desatualizado          string        `json:"desatualizado"`
	StatusEscrita          string        `json:"statusEscrita"`
	Medidor                string        `json:"medidor"`
	Estacao                string        `json:"estacao"`
	Modem                  string        `json:"modem"`
	PermissaoExcluir       bool          `json:"permissaoExcluir"`
	PermissaoList          bool          `json:"permissaoList"`
	AgrVinculoList         []Affiliation `json:"agrVinculoList"`
	CodAgrPai              bool          `json:"codAgrPai"`
}

type PertinentVariables struct {
	CodVar []int
}

type Measurements struct {
	Variables map[int]ReceivedData
}

type ReceivedData struct {
	CodVar int
	Error  error
	Data   Variable
}

type VariableLastData struct {
	ValueType int
	Timestamp string
}
