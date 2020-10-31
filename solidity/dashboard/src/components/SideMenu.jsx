import React, { useState, useContext } from "react"
import { Link, useRouteMatch } from "react-router-dom"
import { Web3Status } from "./Web3Status"
import { Web3Context } from "./WithWeb3Context"
import { ContractsDataContext } from "./ContractsDataContextProvider"
import AddressShortcut from "./AddressShortcut"
import { NetworkStatus } from "./NetworkStatus"
import * as Icons from "./Icons"

export const SideMenuContext = React.createContext({})

export const SideMenuProvider = (props) => {
  const [isOpen, setIsOpen] = useState(false)

  const toggle = () => {
    setIsOpen(!isOpen)
  }

  return (
    <SideMenuContext.Provider value={{ isOpen, toggle }}>
      {props.children}
    </SideMenuContext.Provider>
  )
}

export const SideMenu = (props) => {
  const { isOpen } = useContext(SideMenuContext)
  const { yourAddress, provider } = useContext(Web3Context)
  const { isKeepTokenContractDeployer } = useContext(ContractsDataContext)

  const isDisabled = !yourAddress || !provider

  return (
    <nav
      className={`${isOpen ? "active " : ""}side-menu ${
        isDisabled ? " disabled" : ""
      }`}
    >
      <ul title={isDisabled ? "Please choose a wallet first." : ""}>
        <NavLink
          exact
          to="/liquidations"
          label="liquidations"
          icon={<Icons.Rewards />}
        />
        <Web3Status />
        <div className="account-address">
          <h5 className="text-grey-50">
            <span>address:&nbsp;</span>
            <AddressShortcut classNames="h5" address={yourAddress} />
          </h5>
          <NetworkStatus />
        </div>
      </ul>
    </nav>
  )
}

const NavLink = ({
  label,
  to,
  exact,
  icon,
  sublinks,
  wrapperClassName,
  activeClassName,
  withArrowRight,
}) => {
  const match = useRouteMatch({
    path: to,
    exact,
  })

  return (
    <li className={`${wrapperClassName} ${match ? activeClassName : ""}`}>
      <Link to={to}>
        {icon}
        <span className="ml-1">{label}</span>
        {withArrowRight && <Icons.ArrowRight />}
      </Link>
      <SubNavLinks sublinks={sublinks} />
    </li>
  )
}

NavLink.defaultProps = {
  wrapperClassName: "text-label",
  activeClassName: "active-page-link",
  withArrowRight: true,
}

const SubNavLinks = ({ sublinks }) => {
  if (!sublinks) return null

  return (
    <ul className="sublinks">
      {sublinks.map((sublink) => (
        <NavLink
          key={sublink.label}
          {...sublink}
          wrapperClassName="sublink"
          withArrowRight={false}
        />
      ))}
    </ul>
  )
}
