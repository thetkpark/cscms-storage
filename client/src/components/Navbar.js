import Icon from './Icon'
import { Fragment } from 'react'
import styles from '../styles/Navbar.module.css'
const Navbar = ({ auth, handleAction }) => {
	return (
		<Fragment>
			<nav className={styles.NavBar}>
				{auth ? (
					<Fragment>
						<div onClick={() => handleAction('logout')}>
							<Icon name="logout" /> Logout
						</div>
					</Fragment>
				) : (
					<Fragment>
						<div onClick={() => handleAction('login')}>Login</div>
						<div onClick={() => handleAction('register')}>Register</div>
					</Fragment>
				)}
			</nav>
		</Fragment>
	)
}

export default Navbar
