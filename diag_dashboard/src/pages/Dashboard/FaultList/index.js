import TableHeader from './TableHeader';
import TableContent from './TableContent';

export default function FaultList({data}){
    return (
        <div className="fault-table">
            <div className="fault-table-title">
                车辆故障信息概览
            </div>
            <TableHeader/>
            <TableContent data={data}/>
        </div>
    )
}