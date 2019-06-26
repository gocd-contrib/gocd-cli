package configrepo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/gocd-contrib/gocd-cli/api"
	"github.com/gocd-contrib/gocd-cli/dub"
	"github.com/gocd-contrib/gocd-cli/utils"
)

var Model = &Crud{EtagCache: make(map[string]string)}

type Crud struct {
	EtagCache map[string]string
}

func (c *Crud) FetchRepo(id string, then func(*ConfigRepo) error) error {
	utils.Debug(`Fetching config-repo %q from server`, id)

	if err := checkBlankId(id); err != nil {
		return utils.InspectError(err, `validating config-repo id is not blank`)
	}

	return api.V1.Get(c.url(id)).Send(c.fromJsonAnd(id, then), c.onFail)
}

func (c *Crud) DeleteRepo(id string, then func(msg api.MessageResponse) error) error {
	utils.Debug(`Deleting config-repo %q on server`, id)

	if err := checkBlankId(id); err != nil {
		return utils.InspectError(err, `validating config-repo id is not blank`)
	}

	return api.V1.Delete(c.url(id), nil).Send(c.toMsgAnd(then), c.onFail)
}

func (c *Crud) AddRepo(repo *ConfigRepo, then func(*ConfigRepo) error) error {
	utils.Debug(`Creating config-repo %q on server`, repo.Id)

	if err := checkBlankId(repo.Id); err != nil {
		return utils.InspectError(err, `validating config-repo id is not blank`)
	}

	return api.V1.Post(`/api/admin/config_repos`, api.JsonReader(repo), addContentType).
		Send(c.fromJsonAnd(repo.Id, then), c.onFail)
}

func (c *Crud) UpdateRepo(repo *ConfigRepo, then func(*ConfigRepo) error) error {
	utils.Debug(`Updating config-repo %q on server`, repo.Id)

	if err := checkBlankId(repo.Id); err != nil {
		return utils.InspectError(err, `validating config-repo id is not blank`)
	}

	addEtag := func(req *dub.Request) error {
		if etag, ok := c.EtagCache[repo.Id]; ok {
			req.Header(`If-Match`, etag)
			return nil
		}
		return fmt.Errorf(`ETag is not known or has not been fetched for config-repo %q; please`+
			` GET the current structure of this config-repo from the GoCD server (via API) before`+
			` attempting to update it.`, repo.Id)
	}

	return api.V1.Put(c.url(repo.Id), api.JsonReader(repo), addContentType, addEtag).
		Send(c.fromJsonAnd(repo.Id, then), c.onFail)
}

func (c *Crud) onFail(res *dub.Response) error {
	return api.ReadBodyAndDo(res, func(b []byte) error {
		api.DieOnAuthError(res)

		if res.Raw.Request.Method != `POST` {
			_, id := path.Split(res.Raw.Request.URL.Path)
			api.DieOnNotFound(res, `No such config-repo with id: %q`, id)
			api.DieOnEtagStale(res, `Failed to update config-repo with id %q because it was updated by someone else first`, id)
		}

		api.DieOn4XX(res, b, api.ParseCrMessageWithErrors)
		api.DieOnUnexpected(res, b)

		return nil
	})
}

func (c *Crud) fromJsonAnd(id string, then func(*ConfigRepo) error) func(*dub.Response) error {
	return func(res *dub.Response) error {
		if etag := res.Headers.Get(`ETag`); etag != `` {
			c.EtagCache[id] = etag
		}

		return api.ReadBodyAndDo(res, func(b []byte) error {
			v := &ConfigRepo{Configuration: make([]Property, 0)}
			if err := json.Unmarshal(b, v); err == nil {
				if then == nil {
					return nil
				}
				return utils.InspectError(then(v), `executing config-repo %s success hook`, res.Raw.Request.Method)
			} else {
				return utils.InspectError(err, `parsing config-repo object response`)
			}
		})
	}
}

func (c *Crud) toMsgAnd(then func(api.MessageResponse) error) func(*dub.Response) error {
	return func(res *dub.Response) error {
		return api.ReadBodyAndDo(res, func(b []byte) error {
			if msg, err := api.ParseMessage(b); err != nil {
				return utils.InspectError(err, `parsing config-repo message response: %s`, string(b))
			} else {
				if then != nil {
					return utils.InspectError(then(msg), `executing config-repo message response handler`)
				}
				return nil
			}
		})
	}
}

func (c *Crud) url(id string) string {
	return path.Join(`/api/admin/config_repos`, url.PathEscape(id))
}

func checkBlankId(id string) error {
	if `` == strings.TrimSpace(id) {
		return fmt.Errorf(`config-repo id is missing`)
	}
	return nil
}

func addContentType(req *dub.Request) error {
	req.Header(`Content-Type`, `application/json`)
	return nil
}
