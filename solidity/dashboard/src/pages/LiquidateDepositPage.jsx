import React, { 
  // useMemo,
  useCallback, useState, useEffect } from "react"
import { useParams } from "react-router-dom"
import FallingPriceAuctionTable from "../components/FallingPriceAuctionTable"
import StatusBadge, { BADGE_STATUS } from "../components/StatusBadge"
import { useTokensPageContext } from "../contexts/TokensPageContext"
import {
  // ArbitrageurTokenDetails,
   DepositAuctionOverview
} from "../components/DepositLiquidationOverview"
import moment from "moment"
import { LoadingOverlay } from "../components/Loadable"
import DataTableSkeleton from "../components/skeletons/DataTableSkeleton"
import { ViewAddressInBlockExplorer } from "../components/ViewInBlockExplorer"
// import { SubmitButton } from "../components/Button"
import { useShowMessage, messageType } from "../components/Message"
import { useWeb3Context } from "../components/WithWeb3Context"
import { liquidationService } from "../services/tbtc-liquidation.service"
import { useModal } from "../hooks/useModal"
import { useFetchData } from "../hooks/useFetchData"
import {
  satsToTBtcViaWeitoshi,
  displayAmount
} from "../utils/token.utils"
// import { colors } from "../constants/colors"

const LiquidateDepositPage = (props) => {
  const {
    tokensContext,
    isFetching,
    grantsAreFetching,

  } = useTokensPageContext()


  // const { depositAddress } = props
  const web3Context = useWeb3Context()
  const { depositAddress } = useParams()
  const { openConfirmationModal } = useModal()

  const [lastRefreshedMoment, setLastRefreshedMoment] = useState(moment())

  const [stCurrentAuctionValue, , refreshCurrentAuctionValue] = useFetchData(
    liquidationService.getDepositCurrentAuctionValue,
    {},
    depositAddress
  )
  const {
    isFetching: currentAuctionValueIsFetching,
    data: auctionValueBN,
  } = stCurrentAuctionValue

  const [stUserTBtcBalance, , refreshUserTBtcBalance] = useFetchData(
    liquidationService.getTBtcBalanceOf,
    {})
  const {
    isFetching: userTBtcBalanceIsFetching,
    data: tBtcBalance,
  } = stUserTBtcBalance

  const [stDepositState, , refreshDepositState] = useFetchData(
    liquidationService.getDepositState,
    {},
    depositAddress
    )
  const {
    isFetching: depositStateIsFetching,
    data: depositStateObj,
  } = stDepositState

  const [stAuctionSchedule] = useFetchData(
    liquidationService.getDepositAuctionOfferingSchedule,
    {},
    depositAddress)
  const {
    isFetching: auctionScheduleIsFetching,
    data: auctionOfferingSchedule,
  } = stAuctionSchedule

  const [stLastStartedLiquidationEvent] = useFetchData(
    liquidationService.getLastStartedLiquidationEvent,
    {},
    depositAddress
  )
  const {
    isFetching: lastStartedLiquidationEventIsFetching,
    // data: startedLiquidationEvent,
  } = stLastStartedLiquidationEvent

  const [stDepositBondAmount] = useFetchData(
    liquidationService.getDepositEthBalance,
    {},
    depositAddress
  )
  const {
    isFetching: bondAmountIsFetching,
    data: bondAmountWei,
  } = stDepositBondAmount

  const [stDepositSizeSatoshis] = useFetchData(
    liquidationService.getDepositSizeSatoshis,
    {},
    depositAddress
  )
  const {
    isFetching: depositSizeSatoshisIsFetching,
    data: depositSizeSatoshis,
  } = stDepositSizeSatoshis

  const refreshData = useCallback(
    () => {
    refreshCurrentAuctionValue()
    refreshUserTBtcBalance()
    refreshDepositState()
  },
  [refreshCurrentAuctionValue, refreshUserTBtcBalance, refreshDepositState]
  )
  useEffect(
    () => {
      setTimeout(
        () => {
          refreshData()
          setLastRefreshedMoment(moment())
          // console.log(`I refreshed at ${moment().toString()}`)
        },
        75000
        // 15000
      )
    },
    [lastRefreshedMoment, refreshData]
  )
  // startRefreshDataTimer()

  const confirmationModalOptions = useCallback(() => {
    if (bondAmountIsFetching || depositSizeSatoshisIsFetching)
      return {}
    else {
      return {
        modalOptions: { title: "Purchase ETH Bond" },
        title: "You’re about to purchase ETH with tBTC.",
        subtitle:
          `This transaction will spend ${satsToTBtcViaWeitoshi(depositSizeSatoshis).toString()} tBTC to obtain ${displayAmount(bondAmountWei, false)} ETH (or more, depending on the block it goes through).
           It can fail if deposit state changes before this transaction gets accepted. 
           Transaction can also fail if you don’t have enough tBTC`,
        btnText: "Purchase ETH",
        confirmationText: "Y",
      }
    }
  }, [satsToTBtcViaWeitoshi, displayAmount, bondAmountIsFetching, depositSizeSatoshisIsFetching]) 

  const getPercentageOnOffer = useCallback(() => {
    const utils = web3Context.web3.utils
    let pct = utils.toBN(0)
    if (!currentAuctionValueIsFetching && !bondAmountIsFetching) {
      pct = auctionValueBN.mul(utils.toBN(100))
        .div(utils.toBN(bondAmountWei))
    }
    // console.log(`getPercentageOnOffer: ${pct}`)
    return pct
  }, [web3Context, bondAmountIsFetching, bondAmountWei, currentAuctionValueIsFetching, auctionValueBN])

  const showMessage = useShowMessage()

  // const onLiquidateFromSummaryBtn = async () => {
  //   try {
  //     await liquidationService.depositNotifySignatureTimeout(
  //       web3Context,
  //       depositAddress
  //     )
  //     showMessage({
  //       type: messageType.SUCCESS,
  //       title: "Success",
  //       content: "Top up committed successfully",
  //     })
  //   } catch (error) {
  //     showMessage({
  //       type: messageType.ERROR,
  //       title: "Commit action has failed ",
  //       content: error.message,
  //     })
  //     throw error
  //   }
  // }

    const handleSubmit = async (onTransactionHashCallback) => {
    try {
      await openConfirmationModal(confirmationModalOptions())
      // const depositSizeWeitoshi = fromTokenUnit(depositSizeSatoshis, 10)
      await liquidationService.purchaseDepositAtAuction(
        web3Context,
        depositAddress,
        onTransactionHashCallback
      )
      showMessage({
        type: messageType.SUCCESS,
        title: "Success",
        content: "Staking delegate transaction has been successfully completed",
      })
    } catch (error) {
      console.error(error)
      showMessage({
        type: messageType.ERROR,
        title: "Staking delegate action has failed ",
        content: error.message,
      })
      throw error
    }
  }


  //TODO: Fix first section. Conditional should be at a higher level and check if it is a deposit.
  // If it's not in liquidation (has ever been) it will fail startedLiquidationEvent.returnValues
  return (
    // <PageWrapper title="Liquidations">
    <section>
      <div className="flex wrap self-center mb-2">
        <h2 className="text-grey-70">
          {satsToTBtcViaWeitoshi(depositSizeSatoshis).toString()}{` TBTC Deposit`}
        </h2>
        
        {lastStartedLiquidationEventIsFetching === false && depositStateIsFetching === false && (
          <>
            <span className="flex self-center ml-2">
              <ViewAddressInBlockExplorer address={depositAddress} urlSuffix={""} />
            </span>
            <span className="flex self-center ml-2">
              <StatusBadge
                className="self-center"
                status={BADGE_STATUS.DISABLED}
                text={depositStateObj.name}
              />
              {/* <span className="h4 text-grey-50 ml-1">
                {!lastStartedLiquidationEventIsFetching &&
                  moment.unix(startedLiquidationEvent.returnValues._timestamp).toString()}
              </span> */}
            </span>
          </>
        )}
      </div>
      <>
        <DepositAuctionOverview
          auctionOfferSummaryIsFetching={currentAuctionValueIsFetching || bondAmountIsFetching || userTBtcBalanceIsFetching || depositSizeSatoshisIsFetching || depositStateIsFetching}
          depositState={depositStateObj}
          tBtcBalance={tBtcBalance}
          auctionValueBN={auctionValueBN}
          bondAmountWei={bondAmountWei}
          depositSizeSatoshis={depositSizeSatoshis}
          refreshData={refreshData}
          getPercentageOnOffer={getPercentageOnOffer}
          onLiquidateFromSummaryBtn={handleSubmit}
        >
        </DepositAuctionOverview>
      </>

      <LoadingOverlay
        isFetching={
          (tokensContext === "granted" ? grantsAreFetching : isFetching) || auctionScheduleIsFetching
        }
        skeletonComponent={<DataTableSkeleton />}
      >
        <FallingPriceAuctionTable
          auctionScheduleData={auctionOfferingSchedule}
          // cancelStakeSuccessCallback={cancelStakeSuccessCallback}
        />
      </LoadingOverlay>
    </section>
    // </PageWrapper>
  )
}


export default LiquidateDepositPage
