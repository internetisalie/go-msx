package config

import (
	"context"
	"cto-github.cisco.com/NFV-BU/go-msx/types"
	"os"
	"regexp"
	"strings"
)

type envMatcher interface {
	matches(k, v string, names types.StringSet) bool
}

type environmentMatchers []environmentMatcher

func (e environmentMatchers) matches(key, value string, names types.StringSet) bool {
	for _, em := range e {
		if em.matches(key, value, names) {
			return true
		}
	}
	return false
}

type environmentMatcher struct {
	k *regexp.Regexp
	v *regexp.Regexp
}

func (e environmentMatcher) matches(key, value string, names types.StringSet) bool {
	if e.k != nil {
		matches := e.k.MatchString(key)
		if !matches {
			return false
		}

		if names != nil {
			// Find the name submatch
			submatches := e.k.FindStringSubmatch(key)
			if len(submatches) < 2 {
				return false
			}

			if !names.Contains(submatches[1]) {
				return false
			}
		}
	} else if names != nil {
		return false
	}

	if e.v != nil {
		matches := e.v.MatchString(value)
		if !matches {
			return false
		}
	}
	return true
}

type stringTransformer func(string) string

type environmentTransformer struct {
	m envMatcher
	k stringTransformer
}

func (e environmentTransformer) transform(entry ProviderEntry, names types.StringSet) ProviderEntry {
	if e.m.matches(entry.Name, entry.Value, names) {
		entry.Name = e.k(entry.Name)
	}
	return entry
}

type linkTransformer struct {
	linkMatcher   environmentMatcher
	nameExtractor stringTransformer
	transformer   environmentTransformer
}

func (d linkTransformer) getLinkNames(entries ProviderEntries) types.StringSet {
	var results = make(types.StringSet)
	for _, entry := range entries {
		if d.linkMatcher.matches(entry.Name, entry.Value, nil) {
			results.Add(d.nameExtractor(entry.Name))
		}
	}
	return results
}

func (d linkTransformer) transformAll(entries ProviderEntries) {
	linkNames := d.getLinkNames(entries)

	for i := range entries {
		entries[i] = d.transformer.transform(entries[i], linkNames)
	}

	entries.SortByNormalizedName()
}

func newDockerEnvironmentTransformer() linkTransformer {
	const dockerLinkPortSuffix = "_PORT"
	var dockerLinkPortMatcher = environmentMatcher{k: regexp.MustCompile(`^(.*)_PORT$`), v: regexp.MustCompile(`^(tcp|udp)://[\d.:]+$`)}

	return linkTransformer{
		linkMatcher: dockerLinkPortMatcher,
		nameExtractor: func(key string) string {
			return strings.TrimSuffix(key, dockerLinkPortSuffix)
		},
		transformer: environmentTransformer{
			m: environmentMatchers{
				{k: regexp.MustCompile(`^(.*)_NAME$`)},
				dockerLinkPortMatcher,
				{k: regexp.MustCompile(`^(.*)_PORT_\d+_(TCP|UDP)$`)},
				{k: regexp.MustCompile(`^(.*)_PORT_\d+_(TCP|UDP)_(PROTO|PORT|ADDR)$`)},
			},
			k: func(s string) string {
				return "CONTAINER_LINKS_DOCKER_" + s
			},
		},
	}
}

func newKubernetesEnvironmentTransformer() linkTransformer {
	const kubernetesLinkPortSuffix = "_SERVICE_PORT"
	var kubernetesLinkPortMatcher = environmentMatcher{k: regexp.MustCompile(`.*_SERVICE_PORT`), v: regexp.MustCompile(`^\d+$`)}

	return linkTransformer{
		linkMatcher: kubernetesLinkPortMatcher,
		nameExtractor: func(key string) string {
			return strings.TrimSuffix(key, kubernetesLinkPortSuffix)
		},
		transformer: environmentTransformer{
			m: environmentMatchers{
				{k: regexp.MustCompile(`.*_SERVICE_HOST`), v: regexp.MustCompile(`^[\d.]+$`)},
				kubernetesLinkPortMatcher,
			},
			k: func(s string) string {
				return "CONTAINER_LINKS_KUBERNETES_" + s
			},
		},
	}
}

type EnvironmentProvider struct {
	Describer
	SilentNotifier
	linkTransformers []linkTransformer
}

func (p *EnvironmentProvider) Load(ctx context.Context) (ProviderEntries, error) {
	lines := os.Environ()

	var results = make(ProviderEntries, 0, len(lines))

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		entry := NewEntry(p, parts[0], parts[1])
		if entry.NormalizedName == "" || entry.NormalizedName[0] == '.' {
			continue
		}

		results = append(results, entry)
	}

	return results, nil
}

func NewEnvironmentProvider(name string) *EnvironmentProvider {
	return &EnvironmentProvider{
		Describer: Named{
			name: name,
		},
		linkTransformers: []linkTransformer{
			newDockerEnvironmentTransformer(),
			newKubernetesEnvironmentTransformer(),
		},
	}
}
