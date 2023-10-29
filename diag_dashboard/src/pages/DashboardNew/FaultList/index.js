import TableHeader from './TableHeader';
import TableContent from './TableContent';

export default function FaultList({data}){
    return (
        <div className="fault-table">
            <TableHeader/>
            <TableContent data={data}/>
        </div>
    )
}