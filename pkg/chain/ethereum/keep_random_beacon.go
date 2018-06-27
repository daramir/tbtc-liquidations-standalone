package ethereum

import (
	"bufio"
	"fmt"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/keep-network/keep-core/pkg/chain/gen"
)

// KeepRandomBeacon connection information for interface to the contract.
type KeepRandomBeacon struct {
	caller          *gen.KeepRandomBeaconImplV1Caller
	callerOpts      *bind.CallOpts
	transactor      *gen.KeepRandomBeaconImplV1Transactor
	transactorOpts  *bind.TransactOpts
	contract        *gen.KeepRandomBeaconImplV1
	contractAddress common.Address
}

// NewKeepRandomBeacon creates the necessary connections and configurations for
// accessing the contract.
func newKeepRandomBeacon(pv *ethereumChain) (*KeepRandomBeacon, error) {
	contractAddressHex, exists := pv.config.ContractAddresses["KeepRandomBeaconImplV1"]
	if !exists {
		return nil, fmt.Errorf(
			"no address information for 'KeepRandomBeacon' in configuration",
		)
	}
	contractAddress := common.HexToAddress(contractAddressHex)

	beaconTransactor, err := gen.NewKeepRandomBeaconImplV1Transactor(
		contractAddress,
		pv.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a KeepRelayBeaconTranactor contract: [%v]",
			err,
		)
	}

	file, err := os.Open(pv.config.Account.KeyFile)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open keyfile: %s [%v]",
			pv.config.Account.KeyFile,
			err,
		)
	}

	optsTransactor, err := bind.NewTransactor(
		bufio.NewReader(file),
		pv.config.Account.KeyFilePassword,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read keyfile: %s [%v]",
			pv.config.Account.KeyFile,
			err,
		)
	}

	beaconCaller, err := gen.NewKeepRandomBeaconImplV1Caller(
		contractAddress,
		pv.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate a KeepRelayBeaconCaller contract: [%v]",
			err,
		)
	}

	optsCaller := &bind.CallOpts{
		From: contractAddress,
	}

	randomBeaconContract, err := gen.NewKeepRandomBeaconImplV1(
		contractAddress,
		pv.client,
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to instantiate contract at address: %s [%v]",
			contractAddressHex,
			err,
		)
	}

	return &KeepRandomBeacon{
		transactor:      beaconTransactor,
		transactorOpts:  optsTransactor,
		caller:          beaconCaller,
		callerOpts:      optsCaller,
		contract:        randomBeaconContract,
		contractAddress: contractAddress,
	}, nil
}

// Initialized calls the contract and returns true if the contract has
// had its Initialize method called.
func (krb *KeepRandomBeacon) Initialized() (bool, error) {
	return krb.caller.Initialized(krb.callerOpts)
}

// HasMinimumStake returns true if the specified address has sufficient
// state to participate.
func (krb *KeepRandomBeacon) HasMinimumStake(
	address common.Address,
) (bool, error) {
	return krb.caller.HasMinimumStake(krb.callerOpts, address)
}

// RequestRelayEntry start the process of generating a signature.
func (krb *KeepRandomBeacon) RequestRelayEntry(
	blockReward *big.Int,
	rawseed []byte,
) (*types.Transaction, error) {
	seed := big.NewInt(0).SetBytes(rawseed)
	return krb.transactor.RequestRelayEntry(krb.transactorOpts, blockReward, seed)
}

// SubmitGroupPublicKey upon completion of a sgiagure make the contract
// call to put it on chain.
func (krb *KeepRandomBeacon) SubmitGroupPublicKey(
	groupPublicKey []byte,
	requestID *big.Int,
) (*types.Transaction, error) {
	gpk := byteSliceToSliceOf1Byte(groupPublicKey)
	return krb.transactor.SubmitGroupPublicKey(krb.transactorOpts, gpk, requestID)
}

// relayEntryRequestedFunc type of function called for
// RelayEntryRequested event.
type relayEntryRequestedFunc func(
	requestID *big.Int,
	payment *big.Int,
	blockReward *big.Int,
	seed *big.Int,
	blockNumber *big.Int,
)

// WatchRelayEntryRequested watches for event RelayEntryRequested.
func (krb *KeepRandomBeacon) WatchRelayEntryRequested(
	success relayEntryRequestedFunc,
	fail errorCallback,
) error {
	eventChan := make(chan *gen.KeepRandomBeaconImplV1RelayEntryRequested)
	eventSubscription, err := krb.contract.WatchRelayEntryRequested(nil, eventChan)
	if err != nil {
		return fmt.Errorf("error creating watch for RelayEntryRequested events: [%v]", err)
	}
	go func() {
		for {
			select {
			case event := <-eventChan:
				success(
					event.RequestID,
					event.Payment,
					event.BlockReward,
					event.Seed,
					event.BlockNumber,
				)

			case ee := <-eventSubscription.Err():
				fail(ee)
			}
		}
	}()
	return nil
}

// relayEntryGeneratedFunc type of function called for
// RelayEntryGenerated event.
type relayEntryGeneratedFunc func(
	requestID *big.Int,
	RequestResponse *big.Int,
	RequestGroupID *big.Int,
	PreviousEntry *big.Int,
	blockNumber *big.Int,
)

// WatchRelayEntryGenerated watches for event.
func (krb *KeepRandomBeacon) WatchRelayEntryGenerated(
	success relayEntryGeneratedFunc,
	fail errorCallback,
) error {
	eventChan := make(chan *gen.KeepRandomBeaconImplV1RelayEntryGenerated)
	eventSubscription, err := krb.contract.WatchRelayEntryGenerated(nil, eventChan)
	if err != nil {
		return fmt.Errorf("error creating watch for RelayEntryGenerated event: [%v]", err)
	}
	go func() {
		for {
			select {
			case event := <-eventChan:
				success(
					event.RequestID,
					event.RequestResponse,
					event.RequestGroupID,
					event.PreviousEntry,
					event.BlockNumber,
				)

			case ee := <-eventSubscription.Err():
				fail(ee)
			}
		}
	}()
	return nil
}

// relayResetEventFunc type of function called for ResetEvent event.
type relayResetEventFunc func(
	LastValidRelayEntry *big.Int,
	LastValidRelayTxHash *big.Int,
	LastValidRelayBlock *big.Int,
)

// WatchRelayResetEvent watches for event WatchRelayResetEvent.
func (krb *KeepRandomBeacon) WatchRelayResetEvent(
	success relayResetEventFunc,
	fail errorCallback,
) error {
	eventChan := make(chan *gen.KeepRandomBeaconImplV1RelayResetEvent)
	eventSubscription, err := krb.contract.WatchRelayResetEvent(nil, eventChan)
	if err != nil {
		return fmt.Errorf("error creating watch for RelayResetEvent event: [%v]", err)
	}
	go func() {
		for {
			select {
			case event := <-eventChan:
				success(
					event.LastValidRelayEntry,
					event.LastValidRelayTxHash,
					event.LastValidRelayBlock,
				)

			case ee := <-eventSubscription.Err():
				fail(ee)
			}
		}
	}()
	return nil
}

// submitGroupPublicKeyEventFunc type of function called for
// SubmitGroupPublicKeyEvent event.
type submitGroupPublicKeyEventFunc func(
	GroupPublicKey []byte,
	RequestID *big.Int,
	ActivationBlockHeight *big.Int,
)

// WatchSubmitGroupPublicKeyEvent watches for event SubmitGroupPublicKeyEvent.
func (krb *KeepRandomBeacon) WatchSubmitGroupPublicKeyEvent(
	success submitGroupPublicKeyEventFunc,
	fail errorCallback,
) error {
	eventChan := make(chan *gen.KeepRandomBeaconImplV1SubmitGroupPublicKeyEvent)
	eventSubscription, err := krb.contract.WatchSubmitGroupPublicKeyEvent(
		nil,
		eventChan,
	)
	if err != nil {
		return fmt.Errorf("error creating watch for SubmitGroupPublicKeyEvent event: [%v]", err)
	}
	go func() {
		for {
			select {
			case event := <-eventChan:
				gpk := sliceOf1ByteToByteSlice(event.GroupPublicKey)
				success(gpk, event.RequestID, event.ActivationBlockHeight)

			case ee := <-eventSubscription.Err():
				fail(ee)
			}
		}
	}()
	return nil
}