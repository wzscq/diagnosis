import TableCellItem from "./TableCellItem"

export default function TableHeader(){
    
    const background='#070926';

    return (
        <div className="fault-table-header">
            <TableCellItem col={1} row={1} colSpan={1} rowSpan={1} background={background}>
                车辆管理编码
            </TableCellItem>
            <TableCellItem col={2} row={1} colSpan={1} rowSpan={1} background={background}>
                项目号
            </TableCellItem>
            <TableCellItem col={3} row={1} colSpan={1} rowSpan={1} background={background}>
                试验规范
            </TableCellItem>
            <TableCellItem col={4} row={1} colSpan={1} rowSpan={1} background={background}>
                设备号
            </TableCellItem>
            <TableCellItem col={5} row={1} colSpan={1} rowSpan={1} background={background}>
                采集时间
            </TableCellItem>
            <TableCellItem col={6} row={1} colSpan={1} rowSpan={1} background={background}>
                故障类型
            </TableCellItem>
            <TableCellItem col={7} row={1} colSpan={1} rowSpan={1} background={background}>
                状态
            </TableCellItem>
            <TableCellItem col={8} row={1} colSpan={1} rowSpan={1} background={background}>
                备注
            </TableCellItem>
        </div>
    )
}