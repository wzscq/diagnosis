import styles from "./page.module.css";

export default function StatusCard({title,value}){
    return (
        <div className={styles.wrapper}>
            <div>
                <h2>{title}</h2>
                <br/>
                <h1>{value}</h1>
            </div>
        </div>
    )
}