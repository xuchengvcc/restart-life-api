package config

type AIConfig struct {
	ProviderMap map[string]*Provider
	Providers   []*Provider `mapstructure:"providers"`
}

type Provider struct {
	Name     string   `mapstructure:"name"`
	APIKey   string   `mapstructure:"api_key"`
	ModelIDs []string `mapstructure:"model_ids"`
}

func (p *Provider) GetAPIKey() string {
	if p == nil {
		return ""
	}
	return p.APIKey
}

func (p *Provider) GetModelIDs() []string {
	if p == nil {
		return nil
	}
	return p.ModelIDs
}

func setAIConfigDefault() {

}

func (c *AIConfig) PostInit() {
	c.ProviderMap = make(map[string]*Provider)
	c.buildMap()
}

func (c *AIConfig) buildMap() {
	for _, provider := range c.Providers {
		c.ProviderMap[provider.Name] = provider
	}
}
