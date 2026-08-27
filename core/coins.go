package core

import (
	"fmt"
	"strings"
)

// Coins is an amount of money in copper pieces, the SRD's smallest coin.
// Ten copper make a silver and ten silver make a gold.
type Coins int

// Gp returns an amount in gold pieces.
func Gp(n int) Coins { return Coins(n * 100) }

// Sp returns an amount in silver pieces.
func Sp(n int) Coins { return Coins(n * 10) }

// Cp returns an amount in copper pieces.
func Cp(n int) Coins { return Coins(n) }

// String writes the amount the way the SRD price lists do, largest coin
// first, omitting zero denominations: "15 gp", "1 gp 5 sp", "7 cp".
func (c Coins) String() string {
	if c == 0 {
		return "0 cp"
	}
	gp, rest := int(c)/100, int(c)%100
	sp, cp := rest/10, rest%10
	var parts []string
	if gp != 0 {
		parts = append(parts, fmt.Sprintf("%d gp", gp))
	}
	if sp != 0 {
		parts = append(parts, fmt.Sprintf("%d sp", sp))
	}
	if cp != 0 {
		parts = append(parts, fmt.Sprintf("%d cp", cp))
	}
	return strings.Join(parts, " ")
}
