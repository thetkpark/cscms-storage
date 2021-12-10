import React from 'react'
import ReactDOM from 'react-dom'
import { Provider, lightTheme } from '@adobe/react-spectrum'
import './styles/index.css'
import App from './NewApp'
import reportWebVitals from './reportWebVitals'

ReactDOM.render(
	<React.StrictMode>
		<Provider theme={lightTheme} colorScheme="light">
			<App />
		</Provider>
	</React.StrictMode>,
	document.getElementById('root')
)

reportWebVitals()
