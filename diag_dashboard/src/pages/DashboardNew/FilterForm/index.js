import {useState} from 'react';
import { Button,Select,Space } from 'antd';
import {ReloadOutlined,SearchOutlined} from '@ant-design/icons';
import { useDispatch} from 'react-redux';
import axios from 'axios';

import {getHost,getDashboard} from '../../../api';

const { Option } = Select;

const host=getHost()+process.env.REACT_APP_SERVICE_API_PREFIX;

export default function FilterForm(){
  const dispatch=useDispatch();
  const [filters,setFilters]=useState({});
  const [years,setYears]=useState([]);
  const [projects,setProjects]=useState([]);
  const [types,setTypes]=useState([]);
  const [specs,setSpecs]=useState([]);

  const onYearSearch=(value)=>{
    axios({
      url:host+"/dashboard/years",
      data:filters,
      method:"post",
    }).then(res=>{
      setYears(res.data);
    }).catch(err=>{
      console.log(err);
    });
  }

  const onProjectSearch=(value)=>{
    axios({
      url:host+"/dashboard/projects",
      data:filters,
      method:"post",
    }).then(res=>{
      setProjects(res.data);
    }).catch(err=>{
      console.log(err);
    });
  }

  const onTypeSearch=(value)=>{
    axios({
      url:host+"/dashboard/types",
      data:filters,
      method:"post",
    }).then(res=>{
      setTypes(res.data);
    }).catch(err=>{
      console.log(err);
    });
  }

  const onSpecSearch=(value)=>{
    axios({
      url:host+"/dashboard/specs",
      data:filters,
      method:"post",
    }).then(res=>{
      setSpecs(res.data);
    }).catch(err=>{
      console.log(err);
    });
  }

  const onYearFocus=()=>{
    onYearSearch("");
  }

  const onYearChange=(value)=>{
    setFilters({...filters,year:value});
  }

  const onProjectChange=(value)=>{
    setFilters({...filters,project:value});
  }

  const onProjectFocus=()=>{
    onProjectSearch("");
  }

  const onTypeChange=(value)=>{
    setFilters({...filters,type:value});
  }

  const onTypeFocus=()=>{
    onTypeSearch("");
  }

  const onSpecChange=(value)=>{
    setFilters({...filters,spec:value});
  }

  const onSpecFocus=()=>{
    onSpecSearch("");
  }

  const optionYearControls=years?years.map((item,index)=>{
    return (<Option key={item.year} value={item.year}>{item.year}</Option>);
  }):[];

  const optionProjectControls=projects?projects.map((item,index)=>{
    return (<Option key={item.project} value={item.project}>{item.project}</Option>);
  }):[];

  const optionTypeControls=types?types.map((item,index)=>{
    return (<Option key={item.type} value={item.type}>{item.type}</Option>);
  }):[];

  const optionSpecControls=specs?specs.map((item,index)=>{
    return (<Option key={item.spec} value={item.spec}>{item.spec}</Option>);
  }):[];

  const onReset=()=>{
    setFilters({});
    dispatch(getDashboard({}));
  }

  const onSearch=()=>{
    dispatch(getDashboard(filters));
  }

  return (
    <div className="filter-form">
        <div className='filter-item' style={{gridColumnStart:1,gridColumnEnd:2,gridRowStart:1,gridRowEnd:2}}>
            <div className='filter-item-label' style={{padding:'5px',float:'left',width:'100px',color:'black',textAlign:'right'}}>
              年份:
            </div>
            <div className='filter-item-input' style={{float:'left',width:'calc(100% - 100px)'}}>
              <Select style={{width:'100%'}}
                allowClear
                showSearch
                value={filters?.year}
                onChange={onYearChange}
                onFocus={onYearFocus}
                filterOption={(input, option) =>
                    option.children?.toLowerCase().indexOf(input.toLowerCase()) >= 0||
                    option.value?.toLowerCase().indexOf(input.toLowerCase()) >= 0
                }       
              >
                {optionYearControls}
              </Select>
            </div>
        </div>
        <div className='filter-item' style={{gridColumnStart:2,gridColumnEnd:3,gridRowStart:1,gridRowEnd:2}}>
            <div className='filter-item-label' style={{padding:'5px',float:'left',width:'100px',color:'black',textAlign:'right'}}>
              项目编号:
            </div>
            <div className='filter-item-input' style={{float:'left',width:'calc(100% - 100px)'}}>
              <Select style={{width:'100%'}}
                allowClear
                showSearch
                value={filters?.project}
                onChange={onProjectChange}
                onFocus={onProjectFocus}
                filterOption={(input, option) =>
                    option.children?.toLowerCase().indexOf(input.toLowerCase()) >= 0||
                    option.value?.toLowerCase().indexOf(input.toLowerCase()) >= 0
                }     
              >
                {optionProjectControls}
              </Select>
            </div>
        </div>
        <div className='filter-item' style={{gridColumnStart:3,gridColumnEnd:4,gridRowStart:1,gridRowEnd:2}}>
            <div className='filter-item-label' style={{padding:'5px',float:'left',width:'100px',color:'black',textAlign:'right'}}>
              控制器:
            </div>
            <div className='filter-item-input' style={{float:'left',width:'calc(100% - 100px)'}}>
              <Select style={{width:'100%'}}
                allowClear
                showSearch
                value={filters?.type}
                onChange={onTypeChange}
                onFocus={onTypeFocus}
                filterOption={(input, option) =>
                    option.children?.toLowerCase().indexOf(input.toLowerCase()) >= 0||
                    option.value?.toLowerCase().indexOf(input.toLowerCase()) >= 0
                }     
              >
                {optionTypeControls}
              </Select>
            </div>
        </div>
        <div className='filter-item' style={{gridColumnStart:4,gridColumnEnd:5,gridRowStart:1,gridRowEnd:2}}>
            <div className='filter-item-label' style={{padding:'5px',float:'left',width:'100px',color:'black',textAlign:'right'}}>
              试验规范:
            </div>
            <div className='filter-item-input' style={{float:'left',width:'calc(100% - 100px)'}}>
              <Select style={{width:'100%'}}
                allowClear
                showSearch
                value={filters?.spec}
                onChange={onSpecChange}
                onFocus={onSpecFocus}
                filterOption={(input, option) =>
                    option.children?.toLowerCase().indexOf(input.toLowerCase()) >= 0||
                    option.value?.toLowerCase().indexOf(input.toLowerCase()) >= 0
                }     
              >
                {optionSpecControls}
              </Select>
            </div>
        </div>
        <div className='filter-item' style={{gridColumnStart:4,gridColumnEnd:5,gridRowStart:2,gridRowEnd:3}}>
            <Space style={{float:'right'}}>
              <Button onClick={onReset} type='primary' icon={<ReloadOutlined />}>重置</Button>
              <Button onClick={onSearch} type='primary' icon={<SearchOutlined />}>查询</Button>
            </Space>
        </div>
    </div>
  )
}