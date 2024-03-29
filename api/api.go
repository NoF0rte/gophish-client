package api

import (
	"crypto/tls"
	"fmt"
	"html"
	"regexp"

	"github.com/NoF0rte/gophish-client/api/models"
	"github.com/go-resty/resty/v2"
)

// Client interacts with the GoPhish admin API client
type Client struct {
	client *resty.Client
}

// NewClient creates a new client
func NewClient(url string, apiKey string) *Client {
	client := resty.New().
		SetBaseURL(url).
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
	if apiKey != "" {
		client.SetAuthToken(apiKey)
	}

	return &Client{
		client: client,
	}
}

// NewClient creates a new client and attempts to retrieve the API key using the credentials
func NewClientFromCredentials(url string, username string, password string) (*Client, error) {
	c := NewClient(url, "")

	apiKey, err := c.GetAPIKey(username, password)
	if err != nil {
		return nil, err
	}

	c.client.SetAuthToken(apiKey)

	return c, nil
}

func (c *Client) newRequest(result interface{}) *resty.Request {
	req := c.client.R()
	if result != nil {
		req = req.SetResult(result)
	}
	return req
}

func (c *Client) get(path string, result interface{}) (*resty.Response, interface{}, error) {
	resp, err := c.newRequest(result).Get(path)
	if err != nil {
		return nil, nil, err
	}

	r := resp.Result()

	return resp, r, nil
}

func (c *Client) post(path string, body interface{}, result interface{}) (*resty.Response, interface{}, error) {
	req := c.newRequest(result)
	if body != nil {
		req.SetBody(body)
	}

	resp, err := req.Post(path)
	if err != nil {
		return nil, nil, err
	}

	r := resp.Result()

	return resp, r, nil
}

func (c *Client) put(path string, body interface{}, result interface{}) (*resty.Response, interface{}, error) {
	req := c.newRequest(result)
	if body != nil {
		req.SetBody(body)
	}

	resp, err := req.Put(path)
	if err != nil {
		return nil, nil, err
	}

	r := resp.Result()

	return resp, r, nil
}

func (c *Client) delete(path string, body interface{}, result interface{}) (*resty.Response, interface{}, error) {
	req := c.newRequest(result)
	if body != nil {
		req.SetBody(body)
	}

	resp, err := req.Delete(path)
	if err != nil {
		return nil, nil, err
	}

	r := resp.Result()

	return resp, r, nil
}

// GetAPIKey
func (c *Client) GetAPIKey(username string, password string) (string, error) {
	resp, err := c.client.R().Get("/login")
	if err != nil {
		return "", err
	}

	cookies := resp.Cookies()
	csrfTokenRe := regexp.MustCompile(`name="csrf_token"\s*value="([^"]+)"`)

	body := string(resp.Body())
	matches := csrfTokenRe.FindStringSubmatch(body)
	if len(matches) == 0 {
		return "", fmt.Errorf("error finding csrf_token")
	}

	csrfToken := html.UnescapeString(matches[1])

	resp, err = c.client.R().
		SetCookies(cookies).
		SetFormData(map[string]string{
			"username":   username,
			"password":   password,
			"csrf_token": csrfToken,
		}).
		Post("/login")

	if err != nil {
		return "", err
	}

	if resp.IsError() {
		return "", fmt.Errorf("error: %s", resp.Status())
	}

	resp, err = c.client.R().
		SetCookies(resp.Cookies()).
		Get("/settings")
	if err != nil {
		return "", nil
	}

	body = string(resp.Body())
	apiKeyRe := regexp.MustCompile(`api_key\s*:\s*"([^"]+)"`)

	matches = apiKeyRe.FindStringSubmatch(body)
	if len(matches) == 0 {
		return "", fmt.Errorf("error finding api key")
	}

	return matches[1], nil
}

func (c *Client) GetTemplates() ([]*models.Template, error) {
	var templates []*models.Template
	_, _, err := c.get("/api/templates/", &templates)
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func (c *Client) GetTemplateByID(id int) (*models.Template, error) {
	t := &models.Template{}
	_, _, err := c.get(fmt.Sprintf("/api/templates/%d", id), t)
	if err != nil {
		return nil, err
	}

	if t.ID == 0 {
		return nil, nil
	}

	return t, nil
}

func (c *Client) GetTemplateByName(name string) (*models.Template, error) {
	templates, err := c.GetTemplates()
	if err != nil {
		return nil, err
	}

	for _, t := range templates {
		if t.Name == name {
			return t, nil
		}
	}

	return nil, nil
}

func (c *Client) GetTemplatesByRegex(re string) ([]*models.Template, error) {
	templates, err := c.GetTemplates()
	if err != nil {
		return nil, err
	}

	var filtered []*models.Template
	regex := regexp.MustCompile(re)
	for _, t := range templates {
		if regex.MatchString(t.Name) {
			filtered = append(filtered, t)
		}
	}

	return filtered, nil
}

func (c *Client) GetSendingProfiles() ([]*models.SendingProfile, error) {
	var profiles []*models.SendingProfile
	_, _, err := c.get("/api/smtp/", &profiles)
	if err != nil {
		return nil, err
	}

	return profiles, nil
}

func (c *Client) GetSendingProfileByID(id int) (*models.SendingProfile, error) {
	profile := &models.SendingProfile{}
	_, _, err := c.get(fmt.Sprintf("/api/smtp/%d", id), profile)
	if err != nil {
		return nil, err
	}

	if profile.ID == 0 {
		return nil, nil
	}

	return profile, nil
}

func (c *Client) GetSendingProfileByName(name string) (*models.SendingProfile, error) {
	profiles, err := c.GetSendingProfiles()
	if err != nil {
		return nil, err
	}

	for _, t := range profiles {
		if t.Name == name {
			return t, nil
		}
	}

	return nil, nil
}

func (c *Client) GetSendingProfilesByRegex(re string) ([]*models.SendingProfile, error) {
	profiles, err := c.GetSendingProfiles()
	if err != nil {
		return nil, err
	}

	var filtered []*models.SendingProfile
	regex := regexp.MustCompile(re)
	for _, t := range profiles {
		if regex.MatchString(t.Name) {
			filtered = append(filtered, t)
		}
	}

	return filtered, nil
}

func (c *Client) GetCampaigns() ([]*models.Campaign, error) {
	var campaigns []*models.Campaign
	_, _, err := c.get("/api/campaigns/", &campaigns)
	if err != nil {
		return nil, err
	}

	return campaigns, nil
}

func (c *Client) GetCampaignByID(id int) (*models.Campaign, error) {
	campaign := &models.Campaign{}
	_, _, err := c.get(fmt.Sprintf("/api/campaigns/%d", id), campaign)
	if err != nil {
		return nil, err
	}

	if campaign.ID == 0 {
		return nil, nil
	}

	return campaign, nil
}

func (c *Client) GetCampaignByName(name string) (*models.Campaign, error) {
	campaigns, err := c.GetCampaigns()
	if err != nil {
		return nil, err
	}

	for _, c := range campaigns {
		if c.Name == name {
			return c, nil
		}
	}

	return nil, nil
}

func (c *Client) GetCampaignsByRegex(re string) ([]*models.Campaign, error) {
	campaigns, err := c.GetCampaigns()
	if err != nil {
		return nil, err
	}

	var filtered []*models.Campaign
	regex := regexp.MustCompile(re)
	for _, c := range campaigns {
		if regex.MatchString(c.Name) {
			filtered = append(filtered, c)
		}
	}

	return filtered, nil
}

func (c *Client) GetLandingPages() ([]*models.Page, error) {
	var pages []*models.Page
	_, _, err := c.get("/api/pages/", &pages)
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func (c *Client) GetLandingPageByID(id int) (*models.Page, error) {
	page := &models.Page{}
	_, _, err := c.get(fmt.Sprintf("/api/pages/%d", id), page)
	if err != nil {
		return nil, err
	}

	if page.ID == 0 {
		return nil, nil
	}

	return page, nil
}

func (c *Client) GetLandingPageByName(name string) (*models.Page, error) {
	pages, err := c.GetLandingPages()
	if err != nil {
		return nil, err
	}

	for _, t := range pages {
		if t.Name == name {
			return t, nil
		}
	}

	return nil, nil
}

func (c *Client) GetLandingPagesByRegex(re string) ([]*models.Page, error) {
	pages, err := c.GetLandingPages()
	if err != nil {
		return nil, err
	}

	var filtered []*models.Page
	regex := regexp.MustCompile(re)
	for _, p := range pages {
		if regex.MatchString(p.Name) {
			filtered = append(filtered, p)
		}
	}

	return filtered, nil
}

func (c *Client) GetGroups() ([]*models.Group, error) {
	var groups []*models.Group
	_, _, err := c.get("/api/groups/", &groups)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (c *Client) GetGroupByID(id int) (*models.Group, error) {
	group := &models.Group{}
	_, _, err := c.get(fmt.Sprintf("/api/groups/%d", id), group)
	if err != nil {
		return nil, err
	}

	if group.ID == 0 {
		return nil, nil
	}

	return group, nil
}

func (c *Client) GetGroupByName(name string) (*models.Group, error) {
	groups, err := c.GetGroups()
	if err != nil {
		return nil, err
	}

	for _, g := range groups {
		if g.Name == name {
			return g, nil
		}
	}

	return nil, nil
}

func (c *Client) GetGroupsByRegex(re string) ([]*models.Group, error) {
	groups, err := c.GetGroups()
	if err != nil {
		return nil, err
	}

	var filtered []*models.Group
	regex := regexp.MustCompile(re)
	for _, g := range groups {
		if regex.MatchString(g.Name) {
			filtered = append(filtered, g)
		}
	}

	return filtered, nil
}

func (c *Client) GetGroupsSummary() ([]*models.Group, error) {
	summary := &models.GroupsSummary{}
	_, _, err := c.get("/api/groups/summary", summary)
	if err != nil {
		return nil, err
	}

	return summary.Groups, nil
}

func (c *Client) GetGroupSummaryByID(id int) (*models.Group, error) {
	group := &models.Group{}
	_, _, err := c.get(fmt.Sprintf("/api/groups/%d/summary", id), group)
	if err != nil {
		return nil, err
	}

	if group.ID == 0 {
		return nil, nil
	}

	return group, nil
}

func (c *Client) GetGroupSummaryByName(name string) (*models.Group, error) {
	groups, err := c.GetGroupsSummary()
	if err != nil {
		return nil, err
	}

	for _, g := range groups {
		if g.Name == name {
			return g, nil
		}
	}

	return nil, nil
}

func (c *Client) GetGroupsSummaryByRegex(re string) ([]*models.Group, error) {
	groups, err := c.GetGroupsSummary()
	if err != nil {
		return nil, err
	}

	var filtered []*models.Group
	regex := regexp.MustCompile(re)
	for _, g := range groups {
		if regex.MatchString(g.Name) {
			filtered = append(filtered, g)
		}
	}

	return filtered, nil
}

func (c *Client) DeleteTemplateByID(id int64) (*models.GenericResponse, error) {
	r := &models.GenericResponse{}
	_, _, err := c.delete(fmt.Sprintf("/api/templates/%d", id), nil, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) DeleteTemplateByName(name string) (*models.GenericResponse, error) {
	templates, err := c.GetTemplates()
	if err != nil {
		return nil, err
	}

	var template *models.Template
	for _, t := range templates {
		if t.Name == name {
			template = t
			break
		}
	}

	if template == nil {
		return nil, fmt.Errorf("template %s not found", name)
	}

	return c.DeleteTemplateByID(template.ID)
}

func (c *Client) CreateTemplate(template *models.Template) (*models.Template, error) {
	template.ID = 0 // Ensure the ID is always 0

	_, result, err := c.post("/api/templates/", template, &models.Template{})
	if err != nil {
		return nil, err
	}

	return result.(*models.Template), nil
}

func (c *Client) CreateSendingProfile(profile *models.SendingProfile) (*models.SendingProfile, error) {
	profile.ID = 0 // Ensure the ID is always 0

	if profile.Interface == "" {
		profile.Interface = models.InterfaceSMTP
	}

	_, result, err := c.post("/api/smtp/", profile, &models.SendingProfile{})
	if err != nil {
		return nil, err
	}

	return result.(*models.SendingProfile), nil
}

func (c *Client) CreateLandingPage(page *models.Page) (*models.Page, error) {
	page.ID = 0 // Ensure the ID is always 0

	_, result, err := c.post("/api/pages/", page, &models.Page{})
	if err != nil {
		return nil, err
	}

	return result.(*models.Page), nil
}

func (c *Client) CreateGroup(group *models.Group) (*models.Group, error) {
	group.ID = 0 // Ensure the ID is always 0

	_, result, err := c.post("/api/groups/", group, &models.Group{})
	if err != nil {
		return nil, err
	}

	return result.(*models.Group), nil
}

func (c *Client) CreateCampaign(campaign *models.Campaign) (*models.Campaign, error) {
	newCampaign := &models.Campaign{
		Name:       campaign.Name,
		URL:        campaign.URL,
		LaunchDate: campaign.LaunchDate,
		SendByDate: campaign.SendByDate,
	}

	if campaign.Template != nil {
		newCampaign.Template = &models.Template{
			Name: campaign.Template.Name,
		}
	}

	if campaign.Page != nil {
		newCampaign.Page = &models.Page{
			Name: campaign.Page.Name,
		}
	}

	if campaign.SMTP != nil {
		newCampaign.SMTP = &models.SendingProfile{
			Name: campaign.SMTP.Name,
		}
	}

	if len(campaign.Groups) > 0 {
		var groups []*models.Group
		for _, g := range campaign.Groups {
			groups = append(groups, &models.Group{
				Name: g.Name,
			})
		}
		newCampaign.Groups = groups
	}

	_, result, err := c.post("/api/campaigns/", newCampaign, &models.Campaign{})
	if err != nil {
		return nil, err
	}

	return result.(*models.Campaign), nil
}

func (c *Client) UpdateTemplate(id int64, template *models.Template) (*models.Template, error) {
	template.ID = id
	_, result, err := c.put(fmt.Sprintf("/api/templates/%d", id), template, &models.Template{})
	if err != nil {
		return nil, err
	}

	return result.(*models.Template), nil
}

func (c *Client) UpdateSendingProfile(id int64, profile *models.SendingProfile) (*models.SendingProfile, error) {
	profile.ID = id
	_, result, err := c.put(fmt.Sprintf("/api/smtp/%d", id), profile, &models.SendingProfile{})
	if err != nil {
		return nil, err
	}

	return result.(*models.SendingProfile), nil
}

func (c *Client) UpdateLandingPage(id int64, page *models.Page) (*models.Page, error) {
	page.ID = id
	_, result, err := c.put(fmt.Sprintf("/api/pages/%d", id), page, &models.Page{})
	if err != nil {
		return nil, err
	}

	return result.(*models.Page), nil
}

func (c *Client) UpdateGroup(id int64, group *models.Group) (*models.Group, error) {
	group.ID = id
	_, result, err := c.put(fmt.Sprintf("/api/groups/%d", id), group, &models.Group{})
	if err != nil {
		return nil, err
	}

	return result.(*models.Group), nil
}

func (c *Client) DeleteSendingProfileByID(id int64) (*models.GenericResponse, error) {
	r := &models.GenericResponse{}
	_, _, err := c.delete(fmt.Sprintf("/api/smtp/%d", id), nil, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) DeleteSendingProfileByName(name string) (*models.GenericResponse, error) {
	profile, err := c.GetSendingProfileByName(name)
	if err != nil {
		return nil, err
	}

	if profile == nil {
		return nil, fmt.Errorf("profile %s not found", name)
	}

	return c.DeleteSendingProfileByID(profile.ID)
}

func (c *Client) DeleteLandingPageByID(id int64) (*models.GenericResponse, error) {
	r := &models.GenericResponse{}
	_, _, err := c.delete(fmt.Sprintf("/api/pages/%d", id), nil, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) DeleteLandingPageByName(name string) (*models.GenericResponse, error) {
	page, err := c.GetLandingPageByName(name)
	if err != nil {
		return nil, err
	}

	if page == nil {
		return nil, fmt.Errorf("landing page %s not found", name)
	}

	return c.DeleteLandingPageByID(page.ID)
}

func (c *Client) DeleteGroupByID(id int64) (*models.GenericResponse, error) {
	r := &models.GenericResponse{}
	_, _, err := c.delete(fmt.Sprintf("/api/groups/%d", id), nil, r)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *Client) DeleteGroupByName(name string) (*models.GenericResponse, error) {
	group, err := c.GetGroupByName(name)
	if err != nil {
		return nil, err
	}

	if group == nil {
		return nil, fmt.Errorf("group %s not found", name)
	}

	return c.DeleteGroupByID(group.ID)
}

func (c *Client) ImportSite(req models.ImportSite) (string, error) {
	_, result, err := c.post("/api/import/site", &req, &models.ImportedSite{})
	if err != nil {
		return "", err
	}

	return (result.(*models.ImportedSite)).HTML, nil
}
