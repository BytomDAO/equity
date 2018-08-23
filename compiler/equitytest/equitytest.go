package equitytest

const TrivialLock = `
contract TrivialLock() locks x of y {
  clause trivialUnlock() {
    unlock x of y
  }
}
`

const LockWithPublicKey = `
contract LockWithPublicKey(publicKey: PublicKey) locks x of y {
  clause unlockWithSig(sig: Signature) {
    verify checkTxSig(publicKey, sig)
    unlock x of y
  }
}
`

const LockWithPKHash = `
contract LockWithPublicKeyHash(pubKeyHash: Hash) locks x of y {
  clause spend(pubKey: PublicKey, sig: Signature) {
    verify sha3(pubKey) == pubKeyHash
    verify checkTxSig(pubKey, sig)
    unlock x of y
  }
}
`

const LockWith2of3Keys = `
contract LockWith3Keys(pubkey1, pubkey2, pubkey3: PublicKey) locks x of y {
  clause unlockWith2Sigs(sig1, sig2: Signature) {
    verify checkTxMultiSig([pubkey1, pubkey2, pubkey3], [sig1, sig2])
    unlock x of y
  }
}
`

const LockToOutput = `
contract LockToOutput(address: Program) locks x of y {
  clause relock() {
    lock x of y with address
  }
}
`

const TradeOffer = `
contract TradeOffer(requestedAsset: Asset, requestedAmount: Amount, sellerProgram: Program, sellerKey: PublicKey) locks x of y {
  clause trade() {
    lock requestedAmount of requestedAsset with sellerProgram
    unlock x of y
  }
  clause cancel(sellerSig: Signature) {
    verify checkTxSig(sellerKey, sellerSig)
    unlock x of y
  }
}
`

const EscrowedTransfer = `
contract EscrowedTransfer(agent: PublicKey, sender: Program, recipient: Program) locks x of y {
  clause approve(sig: Signature) {
    verify checkTxSig(agent, sig)
    lock x of y with recipient
  }
  clause reject(sig: Signature) {
    verify checkTxSig(agent, sig)
    lock x of y with sender
  }
}
`

const RevealPreimage = `
contract RevealPreimage(hash: Hash) locks x of y {
  clause reveal(string: String) {
    verify sha3(string) == hash
    unlock x of y
  }
}
`
