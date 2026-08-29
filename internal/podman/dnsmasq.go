package podman

import (
	"fmt"
	"io"
	"strings"
)

// DNSMasqImage is the local tag of the dnsmasq image lerd-dns runs.
const DNSMasqImage = "lerd-dnsmasq:local"

const dnsmasqContainerfile = "FROM docker.io/library/alpine:latest\nRUN apk add --no-cache dnsmasq\n"

// BuildDNSMasqImage builds the dnsmasq image, retrying once with nameservers
// pinned via --dns when the plain build fails. apk resolves from inside the
// build's own network namespace rather than through the host resolver that
// just pulled the base image, and podman's handling of a stub resolver there
// does not hold on every rootless host (#1537). Joining the lerd network is
// not an alternative: a rootless build cannot attach to a named network.
func BuildDNSMasqImage(w io.Writer, nameservers []string) error {
	err := buildDNSMasq(w, nil)
	if err == nil || len(nameservers) == 0 {
		return err
	}
	fmt.Fprintf(w, "\nretrying with the host resolvers (%s)\n", strings.Join(nameservers, ", "))
	return buildDNSMasq(w, nameservers)
}

func buildDNSMasq(w io.Writer, nameservers []string) error {
	args := []string{"build", "-t", DNSMasqImage}
	for _, ns := range nameservers {
		args = append(args, "--dns", ns)
	}
	cmd := Cmd(append(args, "-")...)
	cmd.Stdin = strings.NewReader(dnsmasqContainerfile)
	cmd.Stdout = w
	cmd.Stderr = w
	return cmd.Run()
}
