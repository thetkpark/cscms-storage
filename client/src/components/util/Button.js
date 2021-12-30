import styles from '../../styles/Button.module.css'
const Button = ({ children, color, bgColor, action, style }) => {
	return (
		<button className={styles.Btn} style={{ color: color, background: bgColor, ...style }} onClick={action}>
			{children}
		</button>
	)
}

export default Button
