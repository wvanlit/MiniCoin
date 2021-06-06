/**
 * @Author: Wessel van Lit
 * @Project: minicoin
 * @Date: 06-Jun-2021
 */

package blockchain

func CreateThreeWallets() (*Wallet, *Wallet, *Wallet, error) {
	w1, err1 := GenerateWallet()
	if err1 != nil {
		return nil, nil, nil, err1
	}

	w2, err2 := GenerateWallet()
	if err2 != nil {
		return nil, nil, nil, err2
	}

	w3, err3 := GenerateWallet()
	if err3 != nil {
		return nil, nil, nil, err3
	}

	return w1, w2, w3, nil
}
