package main

import (
	"github.com/google/go-github/github"
	"strconv"
	"log"
	"strings"
)

type WebHookStruct struct {
	Url string
	Events []string
	Enabled string
	SSLCheck bool
	identified bool
	ContentType string // json or form. Default form
}

const hook_ignored = "ignored"

func (h WebHookStruct)Update(hook *github.Hook) (dirty bool) {
	if hook == nil {
		return
	}
	if h.Enabled != hook_ignored {
		if b, err := strconv.ParseBool(h.Enabled) ; err != nil {
			log.Printf("hook `Enabled` has an invalid boolean string representation '%s'. Ignored. hook is enabled.",
				*hook.Name)
		} else if hook.Active && *hook.Active != b {
			dirty = true
			hook.Active = &b
		}
	}

	if v, found := hook.Config["url"]; found {
		if d, ok := v.(string) ; ok && d != h.Url {
			dirty = true
			hook.Config["url"] = h.Url
			log.Printf("Hook '%s' url updated from '%s' to '%s'.", *hook.Name, d, h.Url)
		} else {
		dirty = true
		hook.Config["url"] = h.Url
		log.Printf("Hook '%s' url set to '%s'.", *hook.Name, h.Url)
		}
	} else {
		dirty = true
		hook.Config["url"] = h.Url
		log.Printf("Hook '%s' url set to '%s'.", *hook.Name, h.Url)
	}

	if v, found := hook.Config["insecure_ssl"] ; found {
		if v.(string) == "1" && h.SSLCheck {
			dirty = true
			hook.Config["insecure_ssl"] = "0"
			log.Printf("Hook '%s' SSL check updated from 'false' to 'true'.", *hook.Name)
		} else if v.(string) == "0" && !h.SSLCheck {
			dirty = true
			hook.Config["insecure_ssl"] = "1"
			log.Printf("Hook '%s' SSL check updated from 'true' to 'false'.", *hook.Name)
		}
	} else if !h.SSLCheck {
		hook.Config["insecure_ssl"] = "1"
		log.Printf("Hook '%s' SSL check is false.", *hook.Name)
	}

	if v, found := hook.Config["content_type"]; found {
		if d, ok := v.(string) ; ok && d != h.ContentType {
			dirty = true
			hook.Config["content_type"] = h.ContentType
			log.Printf("Hook '%s' SSL check is false.", *hook.Name)
			log.Printf("Hook '%s' content_type is updated from '%s' to '%s'.", *hook.Name, d, h.ContentType)
		} else {
			hook.Config["content_type"] = h.ContentType
			log.Printf("Hook '%s' content_type is set to '%s'.", *hook.Name, h.ContentType)
		}
	} else {
		hook.Config["content_type"] = h.ContentType
		log.Printf("Hook '%s' content_type is set to '%s'.", *hook.Name, h.ContentType)
	}

	if hook.Events != h.Events {
		dirty = true
		hook.Events = h.Events
		log.Printf("Hook '%s' events are updated from '%s' to '%s'.",
			*hook.Name, strings.Join(hook.Events, ","), strings.Join(h.Events, ","))
	}
	return
}
