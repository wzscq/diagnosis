export default function CarNumber({total}){
    const style={
        with:'100%',
        height:'200px',
        paddingTop:'3px',
    }

    return (
        <div style={style}>
            <div className="dashboard-subtitle">
                诊断报告总数
            </div>
            <div className="dashboard-title">
                {total}
            </div>
        </div>
    )
}