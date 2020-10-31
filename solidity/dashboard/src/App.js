import React from "react"
import Web3ContextProvider from "./components/Web3ContextProvider"
import Footer from "./components/Footer"
import Header from "./components/Header"
import Routing from "./components/Routing"
import ContractsDataContextProvider from "./components/ContractsDataContextProvider"
import { Messages } from "./components/Message"
import { SideMenu, SideMenuProvider } from "./components/SideMenu"
import { HashRouter as Router } from "react-router-dom"
// import IpfsRouter from "ipfs-react-router"

import { Provider } from "react-redux"
import store from "./store"
import { ModalContextProvider } from "./components/Modal"

const App = () => (
  <Provider store={store}>
    <Messages>
      <Web3ContextProvider>
        <ModalContextProvider>
          <ContractsDataContextProvider>
            <SideMenuProvider>
              <Router>
                <main>
                  <Header />
                  <aside>
                    <SideMenu />
                  </aside>
                  <div className="content">
                    <Routing />
                  </div>
                  <Footer />
                </main>
              </Router>
            </SideMenuProvider>
          </ContractsDataContextProvider>
        </ModalContextProvider>
      </Web3ContextProvider>
    </Messages>
  </Provider>
)

export default App
