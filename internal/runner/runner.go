package runner

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/projectdiscovery/collaborator"
	colbiid "github.com/projectdiscovery/collaborator/biid"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/notify"
)

const (
	defaultHTTPMessage = "The collaborator server received an {{protocol}} request from {{from}} at {{time}}:\n```\n{{request}}\n{{response}}```"
	defaultDNSMessage  = "The collaborator server received a DNS lookup of type {{type}} for the domain name {{domain}} from {{from}} at {{time}}:\n```{{request}}```"
)

type Runner struct {
	options    *Options
	burpcollab *collaborator.BurpCollaborator
	notifier   *notify.Notify
}

func NewRunner(options *Options) (*Runner, error) {
	burpcollab := collaborator.NewBurpCollaborator()

	notifier, err := notify.NewWithOptions(&notify.Options{
		SlackWebHookUrl:         options.SlackWebHookUrl,
		SlackUsername:           options.SlackUsername,
		SlackChannel:            options.SlackChannel,
		Slack:                   options.Slack,
		DiscordWebHookUrl:       options.DiscordWebHookUrl,
		DiscordWebHookUsername:  options.DiscordWebHookUsername,
		DiscordWebHookAvatarUrl: options.DiscordWebHookAvatarUrl,
		Discord:                 options.Discord,
	})
	if err != nil {
		return nil, err
	}

	return &Runner{options: options, burpcollab: burpcollab, notifier: notifier}, nil
}

func (r *Runner) Run() error {
	// If stdin is present pass everything to webhooks and exit
	if hasStdin() {
		br := bufio.NewScanner(os.Stdin)
		for br.Scan() {
			msg := br.Text()
			gologger.Printf(msg)
			r.notifier.SendNotification(msg)
		}
		os.Exit(0)
	}

	// otherwise works as long term collaborator poll and notify via webhook
	// If BIID passed via cli
	if r.options.BIID != "" {
		gologger.Printf("Using BIID: %s", r.options.BIID)
		r.burpcollab.AddBIID(r.options.BIID)
	} else if r.options.InterceptBIID {
		if os.Getuid() != 0 {
			gologger.Warningf("Command may fail as the program is not running as root and unable to access raw sockets")
		}
		gologger.Printf("Attempting to intercept BIID")
		// otherwise attempt to retrieve it
		biid, err := colbiid.Intercept(time.Duration(r.options.InterceptBIIDTimeout) * time.Second)
		if err != nil {
			return err
		}
		gologger.Printf("BIID found, using: %s", biid)
		r.options.BIID = biid
		r.burpcollab.AddBIID(biid)
	}

	if r.options.BIID == "" {
		return fmt.Errorf("BIID not specified or not found")
	}

	err := r.burpcollab.Poll()
	if err != nil {
		return err
	}

	pollTime := time.Duration(r.options.Interval) * time.Second
	for {
		time.Sleep(pollTime)
		r.burpcollab.Poll()

		for _, httpresp := range r.burpcollab.RespBuffer {
			for _, resp := range httpresp.Responses {
				var at int64
				at, _ = strconv.ParseInt(resp.Time, 10, 64)
				atTime := time.Unix(0, at*int64(time.Millisecond))
				if resp.Protocol == "http" || resp.Protocol == "https" {
					rr := strings.NewReplacer(
						"{{protocol}}", strings.ToUpper(resp.Protocol),
						"{{from}}", resp.Client,
						"{{time}}", atTime.String(),
						"{{request}}", resp.Data.RequestDecoded,
						"{{response}}", resp.Data.ResponseDecoded,
					)

					msg := rr.Replace(r.options.HTTPMessage)
					gologger.Printf(msg)
					r.notifier.SendNotification(msg)
				}
				if resp.Protocol == "dns" {
					rr := strings.NewReplacer(
						"{{type}}", resp.Data.RequestType,
						"{{domain}} ", resp.Data.SubDomain,
						"{{from}}", resp.Client,
						"{{time}}", atTime.String(),
						"{{request}}", resp.Data.RawRequestDecoded,
					)
					msg := rr.Replace(r.options.DNSMessage)
					gologger.Printf(msg)
					r.notifier.SendNotification(msg)
				}
			}
		}

		r.burpcollab.Empty()
	}
}

func (r *Runner) Close() {
	r.burpcollab.Empty()
}
