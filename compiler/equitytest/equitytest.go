package equitytest

const TrivialLock = `
contract TrivialLock() locks amount of asset {
  clause trivialUnlock() {
    unlock amount of asset
  }
}
`

const LockWithPublicKey = `
contract LockWithPublicKey(publicKey: PublicKey) locks amount of asset {
  clause unlockWithSig(sig: Signature) {
    verify checkTxSig(publicKey, sig)
    unlock amount of asset
  }
}
`

const LockWithPKHash = `
contract LockWithPublicKeyHash(pubKeyHash: Hash) locks amount of asset {
  clause spend(pubKey: PublicKey, sig: Signature) {
    verify sha3(pubKey) == pubKeyHash
    verify checkTxSig(pubKey, sig)
    unlock amount of asset
  }
}
`

const LockWith2of3Keys = `
contract LockWith3Keys(pubkey1, pubkey2, pubkey3: PublicKey) locks amount of asset {
  clause unlockWith2Sigs(sig1, sig2: Signature) {
    verify checkTxMultiSig([pubkey1, pubkey2, pubkey3], [sig1, sig2])
    unlock amount of asset
  }
}
`

const LockToOutput = `
contract LockToOutput(address: Program) locks amount of asset {
  clause relock() {
    lock amount of asset with address
  }
}
`

const TradeOffer = `
contract TradeOffer(requestedAsset: Asset, requestedAmount: Amount, sellerProgram: Program, sellerKey: PublicKey) locks amount of asset {
  clause trade() {
    lock requestedAmount of requestedAsset with sellerProgram
    unlock amount of asset
  }
  clause cancel(sellerSig: Signature) {
    verify checkTxSig(sellerKey, sellerSig)
    unlock amount of asset
  }
}
`

const EscrowedTransfer = `
contract EscrowedTransfer(agent: PublicKey, sender: Program, recipient: Program) locks amount of asset {
  clause approve(sig: Signature) {
    verify checkTxSig(agent, sig)
    lock amount of asset with recipient
  }
  clause reject(sig: Signature) {
    verify checkTxSig(agent, sig)
    lock amount of asset with sender
  }
}
`

const RevealPreimage = `
contract RevealPreimage(hash: Hash) locks amount of asset {
  clause reveal(string: String) {
    verify sha3(string) == hash
    unlock amount of asset
  }
}
`
