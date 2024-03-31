'use client'
import { Button,Space } from "antd";
import {DatabaseFilled,ProfileFilled} from '@ant-design/icons';
import { useRouter } from 'next/navigation'

export default function Toolbar({direction,styleid}){
    const router = useRouter()

    const switchDriection=()=>{
        if(direction!=='h'){
            router.push(`/drivingstyle/${styleid}/h`)
        } else {
            router.push(`/drivingstyle/${styleid}/v`)
        }
    }

    const close=()=>{
        const message={
            type:'doOperation',
            data:{
                operationItem:{
                    type:'close',
                    params:{
                        location:'modal'
                    }
                }
            }
        };
        window.parent.postMessage(message,'*');
    }

    return (
        <Space size={2}>
            <Button onClick={switchDriection} disabled={direction==='h'}><ProfileFilled /></Button>,
            <Button onClick={switchDriection} disabled={direction!=='h'}><DatabaseFilled /></Button>,
            <Button onClick={close} type='primary'>关闭</Button>
        </Space>
    )
}