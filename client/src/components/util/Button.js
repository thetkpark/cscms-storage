import styles from '../../styles/Button.module.css'
const Button = ({ children, color, bgColor, action }) => {
	return (
		<button className={styles.Btn} style={{ color: color, background: bgColor }} onClick={action}>
			{children}
		</button>
	)
}

export default Button
