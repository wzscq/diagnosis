import {Table,ConfigProvider} from 'antd';
import zh_CN from 'antd/lib/locale/zh_CN';

const columns = [
  {
    title: '车辆编码',
    dataIndex: 'vehicle_management_code',
    key: 'vehicle_management_code'
  },
  {
    title: '项目编号',
    dataIndex: 'project_num',
    key: 'project_num',
  },
  {
    title: '试验规范',
    dataIndex: 'specifications',
    key: 'specifications',
  },
  {
    title: '设备编号',
    key: 'device_number',
    dataIndex: 'device_number'
  },
  {
    title: '采集时间',
    key: 'time',
    dataIndex: 'time'
  },
  {
    title: '故障控制器',
    key: 'type',
    dataIndex: 'type'
  },
  {
    title: '状态',
    key: 'status',
    dataIndex: 'status',
    render: (text, record) => (
      <span>
        {text==='0'?(<span style={{color:'red'}}>Open</span>):(
            <span style={{color:'#73c0de'}}>Closed</span>
        )}
      </span>
    ),
  },
  {
    title: '备注',
    key: 'remark',
    dataIndex: 'remark'
  }
];

export default function FaultListNew({data}){
  return (
    <ConfigProvider locale={zh_CN}>
    <Table size='small' locale={zh_CN} columns={columns} dataSource={data} />
    </ConfigProvider>
    );
}