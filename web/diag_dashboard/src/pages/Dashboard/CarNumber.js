export default function CarNumber({total}){
    return (
        <>
            <div className="dashboard-subtitle">
                工程车辆总数（辆）
            </div>
            <div className="dashboard-title">
                {total}
            </div>
        </>
    )
}