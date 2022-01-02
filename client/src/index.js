import React from 'react'
import ReactDOM from 'react-dom'
import './styles/index.css'
import App from './NewApp'
import reportWebVitals from './reportWebVitals'
import { RecoilRoot } from 'recoil'
ReactDOM.render(
	<React.StrictMode>
		<RecoilRoot>
			<App />
		</RecoilRoot>
	</React.StrictMode>,
	document.getElementById('root')
)

reportWebVitals()
