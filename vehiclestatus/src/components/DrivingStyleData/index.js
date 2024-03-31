import styles from './page.module.css';

export default function DrivingStyleData({data}){
    return (
        <div className={styles.dataGrid}>
            <div className={styles.title}>项目</div>
            <div className={styles.title}>数值</div>
            <div className={styles.title}>单位</div>
            <div>最高车速</div><div>{data?.speed_max??''}</div><div>km/h</div>
            <div>最大油门踏板深度比例</div><div>{data?.mtpd_ratio??''}</div><div>%</div>
            <div>最大油门变化率</div><div>{data?.mtc_ratio??''}</div><div>%/s</div>
            <div>最大刹车踏板深度比例</div><div>{data?.mbpd_ratio??''}</div><div>%</div>
            <div>最大刹车变化率</div><div>{data?.mbc_ratio??''}</div><div>%/s</div>
            <div>最大方向盘角速度</div><div>{data?.mtwav??''}</div><div>°/s</div>
            <div>最大侧向加速度</div><div>{data?.lateral_acc_max??''}</div><div>m2/s</div>
            <div>最大纵向加速度</div><div>{data?.long_acc_max??''}</div><div>m2/s</div>
        </div>
    )
}