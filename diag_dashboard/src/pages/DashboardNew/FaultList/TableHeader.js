import TableCellItem from "./TableCellItem"

export default function TableHeader(){
    
    const background='#070926';

    return (
        <div className="fault-table-header">
            <TableCellItem color={"#FFF"} col={1} row={1} colSpan={1} rowSpan={1} background={background}>
                车辆编码
            </TableCellItem>
            <TableCellItem color={"#FFF"} col={2} row={1} colSpan={1} rowSpan={1} background={background}>
                项目编号
            </TableCellItem>
            <TableCellItem color={"#FFF"} col={3} row={1} colSpan={1} rowSpan={1} background={background}>
                试验规范
            </TableCellItem>
            <TableCellItem color={"#FFF"} col={4} row={1} colSpan={1} rowSpan={1} background={background}>
                设备编号
            </TableCellItem>
            <TableCellItem color={"#FFF"} col={5} row={1} colSpan={1} rowSpan={1} background={background}>
                采集时间
            </TableCellItem>
            <TableCellItem color={"#FFF"} col={6} row={1} colSpan={1} rowSpan={1} background={background}>
                故障控制器
            </TableCellItem>
            <TableCellItem color={"#FFF"} col={7} row={1} colSpan={1} rowSpan={1} background={background}>
                状态
            </TableCellItem>
            <TableCellItem color={"#FFF"} col={8} row={1} colSpan={1} rowSpan={1} background={background}>
                备注
            </TableCellItem>
        </div>
    )
}