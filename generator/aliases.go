package main

import "strings"

var aliases = map[string]linkableAlias{} // Map of link to link

func initAliases() {
	newLinkableAlias("IaC", smIAC) // TODO: ADD TAGS?
	newLinkableAlias("CiCd", smCiCd)
	newLinkableAlias("CI/CD", smCiCd)                                // TODO: NOT WORKING
	newLinkableAlias("k8s", techs["Kubernetes"])                     // TODO: ok?
	newLinkableAlias("KMS", cloudServices["Key Management Service"]) // TODO: ok?
	newLinkableAlias("ASM", langs["Assembly"])                       // TODO: ok?
	newLinkableAlias("WebAssembly", langs["WASM"])                   // TODO: ok?
}
func newLinkableAlias(name string, linkTo Linkable) {
	out := linkableAlias{name: name, finalLink: linkTo}
	aliases[name] = out
	addLinkable(strings.ToLower(name), out)
}

type linkableAlias struct {
	name      string
	finalLink Linkable
}

func (pg linkableAlias) Title() string {
	return pg.name
}

func (pg linkableAlias) Dst() string { // TODO: ptr?
	return pg.finalLink.Dst()
}
func (pg linkableAlias) EntryType() string {
	return pg.name
}
