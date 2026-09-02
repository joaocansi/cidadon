package entity

import "strings"

// Party is an official Brazilian political party published by the Câmara dos
// Deputados Open Data API. The catalogue is versioned with the application so
// accepting a registration never depends on an external service being online.
type Party struct {
	ID        uint   `json:"id"`
	Acronym   string `json:"sigla"`
	Name      string `json:"nome"`
	SourceURI string `json:"uri"`
}

var officialParties = []Party{
	{36898, "AVANTE", "Avante", "https://dadosabertos.camara.leg.br/api/v2/partidos/36898"},
	{37905, "CIDADANIA", "Cidadania", "https://dadosabertos.camara.leg.br/api/v2/partidos/37905"},
	{37902, "DC", "Democracia Cristã", "https://dadosabertos.camara.leg.br/api/v2/partidos/37902"},
	{36899, "MDB", "Movimento Democrático Brasileiro", "https://dadosabertos.camara.leg.br/api/v2/partidos/36899"},
	{38011, "MISSÃO", "Partido Missão", "https://dadosabertos.camara.leg.br/api/v2/partidos/38011"},
	{37901, "NOVO", "Partido Novo", "https://dadosabertos.camara.leg.br/api/v2/partidos/37901"},
	{36779, "PCdoB", "Partido Comunista do Brasil", "https://dadosabertos.camara.leg.br/api/v2/partidos/36779"},
	{36786, "PDT", "Partido Democrático Trabalhista", "https://dadosabertos.camara.leg.br/api/v2/partidos/36786"},
	{37906, "PL", "Partido Liberal", "https://dadosabertos.camara.leg.br/api/v2/partidos/37906"},
	{36896, "PODE", "Podemos", "https://dadosabertos.camara.leg.br/api/v2/partidos/36896"},
	{37903, "PP", "Progressistas", "https://dadosabertos.camara.leg.br/api/v2/partidos/37903"},
	{38010, "PRD", "Partido Renovação Democrática", "https://dadosabertos.camara.leg.br/api/v2/partidos/38010"},
	{36832, "PSB", "Partido Socialista Brasileiro", "https://dadosabertos.camara.leg.br/api/v2/partidos/36832"},
	{36834, "PSD", "Partido Social Democrático", "https://dadosabertos.camara.leg.br/api/v2/partidos/36834"},
	{36835, "PSDB", "Partido da Social Democracia Brasileira", "https://dadosabertos.camara.leg.br/api/v2/partidos/36835"},
	{36839, "PSOL", "Partido Socialismo e Liberdade", "https://dadosabertos.camara.leg.br/api/v2/partidos/36839"},
	{36844, "PT", "Partido dos Trabalhadores", "https://dadosabertos.camara.leg.br/api/v2/partidos/36844"},
	{36851, "PV", "Partido Verde", "https://dadosabertos.camara.leg.br/api/v2/partidos/36851"},
	{36886, "REDE", "Rede Sustentabilidade", "https://dadosabertos.camara.leg.br/api/v2/partidos/36886"},
	{37908, "REPUBLICANOS", "Republicanos", "https://dadosabertos.camara.leg.br/api/v2/partidos/37908"},
	{37904, "SOLIDARIEDADE", "Solidariedade", "https://dadosabertos.camara.leg.br/api/v2/partidos/37904"},
	{38009, "UNIÃO", "União Brasil", "https://dadosabertos.camara.leg.br/api/v2/partidos/38009"},
}

func OfficialParties() []Party {
	return append([]Party(nil), officialParties...)
}

// OfficialParty returns the canonical party acronym. Inputs are deliberately
// limited to an acronym, avoiding ambiguous free-form party names.
func OfficialParty(value string) (Party, bool) {
	needle := strings.ToUpper(strings.TrimSpace(value))
	for _, party := range officialParties {
		if strings.ToUpper(party.Acronym) == needle {
			return party, true
		}
	}
	return Party{}, false
}
