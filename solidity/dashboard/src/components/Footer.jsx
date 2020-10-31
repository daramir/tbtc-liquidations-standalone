import React from "react"
import * as Icons from "./Icons"

const socialMedia = [
  // { label: "Twitter", url: "https://twitter.com/keep_project" },
  // { label: "Telegram", url: "https://t.me/KeepNetworkOfficial" },
  // { label: "Reddit", url: "https://www.reddit.com/r/KeepNetwork" },
  { label: "Discord", url: "https://chat.keep.network/" },
];

const aboutUs = [
  { label: "Whitepaper", url: "https://keep.network/whitepaper" },
  { label: "Docs", url: "https://docs.keep.network/tbtc/index.html#liquidation"},
  { label: "More Docs", url: "https://bisontrails.co/keep-active-participation/"}
]

const Footer = () => {
  return (
    <footer>
      <Icons.KeepCircle color="#F2F2F2" />
      <ul>{aboutUs.map(renderFooterLinkItem)}</ul>
      <ul>{socialMedia.map(renderFooterLinkItem)}</ul>
      <div className="signature text-smaller text-grey-70">
        <div>(a fork of)</div>
        <div>A Thesis* Build</div>
        <div>© 2020 Keep, SEZC. All Rights Reserved.</div>
      </div>
    </footer>
  )
}

const FooterLinkItem = ({ label, url }) => (
  <li>
    <a href={url} rel="noopener noreferrer" target="_blank">
      <h5>{label}</h5>
    </a>
  </li>
)

const renderFooterLinkItem = (item) => (
  <FooterLinkItem key={item.label} {...item} />
)

export default Footer
