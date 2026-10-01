package logger

import (
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

// Accept returns the provided filter as is.
func Accept(filter Filter) Filter { return filter }

// Ignore returns the provided filter as is.
func Ignore(filter Filter) Filter { return filter }

// AcceptMethod returns a filter that accepts requests with the specified methods.
func AcceptMethod(methods ...string) Filter {
	return func(c *gin.Context) bool {
		reqMethod := strings.ToLower(c.Request.Method)

		for _, method := range methods {
			if strings.ToLower(method) == reqMethod {
				return true
			}
		}

		return false
	}
}

// IgnoreMethod returns a filter that ignores requests with the specified methods.
func IgnoreMethod(methods ...string) Filter {
	return func(c *gin.Context) bool {
		reqMethod := strings.ToLower(c.Request.Method)

		for _, method := range methods {
			if strings.ToLower(method) == reqMethod {
				return false
			}
		}

		return true
	}
}

// AcceptStatus returns a filter that accepts requests with the specified statuses.
func AcceptStatus(statuses ...int) Filter {
	return func(c *gin.Context) bool {
		for _, status := range statuses {
			if status == c.Writer.Status() {
				return true
			}
		}

		return false
	}
}

// IgnoreStatus returns a filter that ignores requests with the specified statuses.
func IgnoreStatus(statuses ...int) Filter {
	return func(c *gin.Context) bool {
		for _, status := range statuses {
			if status == c.Writer.Status() {
				return false
			}
		}

		return true
	}
}

// AcceptStatusGreaterThan returns a filter that accepts requests with statuses greater than the specified status.
func AcceptStatusGreaterThan(status int) Filter {
	return func(c *gin.Context) bool {
		return c.Writer.Status() > status
	}
}

// IgnoreStatusLessThan returns a filter that ignores requests with statuses less than the specified status.
func IgnoreStatusLessThan(status int) Filter {
	return func(c *gin.Context) bool {
		return c.Writer.Status() < status
	}
}

// AcceptStatusGreaterThanOrEqual returns a filter that accepts requests with statuses greater than or equal to the specified status.
func AcceptStatusGreaterThanOrEqual(status int) Filter {
	return func(c *gin.Context) bool {
		return c.Writer.Status() >= status
	}
}

// IgnoreStatusLessThanOrEqual returns a filter that ignores requests with statuses less than or equal to the specified status.
func IgnoreStatusLessThanOrEqual(status int) Filter {
	return func(c *gin.Context) bool {
		return c.Writer.Status() <= status
	}
}

// AcceptPath returns a filter that accepts requests with the specified paths.
func AcceptPath(urls ...string) Filter {
	return func(c *gin.Context) bool {
		for _, url := range urls {
			if c.Request.URL.Path == url {
				return true
			}
		}

		return false
	}
}

// IgnorePath returns a filter that ignores requests with the specified paths.
func IgnorePath(urls ...string) Filter {
	return func(c *gin.Context) bool {
		for _, url := range urls {
			if c.Request.URL.Path == url {
				return false
			}
		}

		return true
	}
}

// AcceptPathContains returns a filter that accepts requests with paths containing the specified parts.
func AcceptPathContains(parts ...string) Filter {
	return func(c *gin.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request.URL.Path, part) {
				return true
			}
		}

		return false
	}
}

// IgnorePathContains returns a filter that ignores requests with paths containing the specified parts.
func IgnorePathContains(parts ...string) Filter {
	return func(c *gin.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request.URL.Path, part) {
				return false
			}
		}

		return true
	}
}

// AcceptPathPrefix returns a filter that accepts requests with paths having the specified prefixes.
func AcceptPathPrefix(prefixs ...string) Filter {
	return func(c *gin.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnorePathPrefix returns a filter that ignores requests with paths having the specified prefixes.
func IgnorePathPrefix(prefixs ...string) Filter {
	return func(c *gin.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				return false
			}
		}

		return true
	}
}

// AcceptPathSuffix returns a filter that accepts requests with paths having the specified suffixes.
func AcceptPathSuffix(prefixs ...string) Filter {
	return func(c *gin.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnorePathSuffix returns a filter that ignores requests with paths having the specified suffixes.
func IgnorePathSuffix(suffixs ...string) Filter {
	return func(c *gin.Context) bool {
		for _, suffix := range suffixs {
			if strings.HasSuffix(c.Request.URL.Path, suffix) {
				return false
			}
		}

		return true
	}
}

// AcceptPathMatch returns a filter that accepts requests with paths matching the specified regular expressions.
func AcceptPathMatch(regs ...regexp.Regexp) Filter {
	return func(c *gin.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request.URL.Path)) {
				return true
			}
		}

		return false
	}
}

// IgnorePathMatch returns a filter that ignores requests with paths matching the specified regular expressions.
func IgnorePathMatch(regs ...regexp.Regexp) Filter {
	return func(c *gin.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request.URL.Path)) {
				return false
			}
		}

		return true
	}
}

// AcceptHost returns a filter that accepts requests with the specified hosts.
func AcceptHost(hosts ...string) Filter {
	return func(c *gin.Context) bool {
		for _, host := range hosts {
			if c.Request.URL.Host == host {
				return true
			}
		}

		return false
	}
}

// IgnoreHost returns a filter that ignores requests with the specified hosts.
func IgnoreHost(hosts ...string) Filter {
	return func(c *gin.Context) bool {
		for _, host := range hosts {
			if c.Request.URL.Host == host {
				return false
			}
		}

		return true
	}
}

// AcceptHostContains returns a filter that accepts requests with hosts containing the specified parts.
func AcceptHostContains(parts ...string) Filter {
	return func(c *gin.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request.URL.Host, part) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostContains returns a filter that ignores requests with hosts containing the specified parts.
func IgnoreHostContains(parts ...string) Filter {
	return func(c *gin.Context) bool {
		for _, part := range parts {
			if strings.Contains(c.Request.URL.Host, part) {
				return false
			}
		}

		return true
	}
}

// AcceptHostPrefix returns a filter that accepts requests with hosts having the specified prefixes.
func AcceptHostPrefix(prefixs ...string) Filter {
	return func(c *gin.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request.URL.Host, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostPrefix returns a filter that ignores requests with hosts having the specified prefixes.
func IgnoreHostPrefix(prefixs ...string) Filter {
	return func(c *gin.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request.URL.Host, prefix) {
				return false
			}
		}

		return true
	}
}

// AcceptHostSuffix returns a filter that accepts requests with hosts having the specified suffixes.
func AcceptHostSuffix(prefixs ...string) Filter {
	return func(c *gin.Context) bool {
		for _, prefix := range prefixs {
			if strings.HasPrefix(c.Request.URL.Host, prefix) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostSuffix returns a filter that ignores requests with hosts having the specified suffixes.
func IgnoreHostSuffix(suffixs ...string) Filter {
	return func(c *gin.Context) bool {
		for _, suffix := range suffixs {
			if strings.HasSuffix(c.Request.URL.Host, suffix) {
				return false
			}
		}

		return true
	}
}

// AcceptHostMatch returns a filter that accepts requests with hosts matching the specified regular expressions.
func AcceptHostMatch(regs ...regexp.Regexp) Filter {
	return func(c *gin.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request.URL.Host)) {
				return true
			}
		}

		return false
	}
}

// IgnoreHostMatch returns a filter that ignores requests with hosts matching the specified regular expressions.
func IgnoreHostMatch(regs ...regexp.Regexp) Filter {
	return func(c *gin.Context) bool {
		for _, reg := range regs {
			if reg.Match([]byte(c.Request.URL.Host)) {
				return false
			}
		}

		return true
	}
}
