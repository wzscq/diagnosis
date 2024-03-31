export async function getVechicleStatus(pageSize,filter) {
    let headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("token","carapiv2");
  
    let body ={
      modelID:"diag_vehicle_status",
          fields:[
        {field:'id'},
        {field:"vehicle_code"},
        {field:"project_code"},
        {field:"test_spec"},
        {field:"device_code"},
        {field:"speed_max"},
        {field:"speed_avg"},
        {field:"mileage"},
        {field:"ec_avg"},
        {field:"travel_time"},
        {field:"create_time"},
        {field:"create_user"},
        {field:"update_time"},
        {field:"update_user"},
        {field:"version"}
      ],
      sorter:[{field:"update_time",order:"desc"}],
      pagination:{current:1,pageSize:pageSize??1000},
      filter:filter
    }
  
    let options = {headers,body:JSON.stringify(body),method:"POST",cache:'no-store'}
    const res = await fetch(`${process.env.CRV_SERVICE_URL}/data/query`,options)
    // The return value is *not* serialized
    // You can return Date, Map, Set, etc.
   
    if (!res.ok) {
      // This will activate the closest `error.js` Error Boundary
      throw new Error('Failed to fetch data')
    }
   
    return res.json()
  }

  export async function getDrivingStyle(id) {
    let headers = new Headers();
    headers.append("Content-Type", "application/json");
    headers.append("token","carapiv2");
  
    let body ={
      modelID:"diag_driving_style",
          fields:[
            {"field":"id"},
            {"field":"vehicle_code"},
            {"field":"project_code"},
            {"field":"test_spec"},
            {"field":"device_code"},
            {"field":"driving_style"},
            {"field":"acquisition_time"},
            {"field":"speed_max"},
            {"field":"mtpd_ratio"},
            {"field":"mtc_ratio"},
            {"field":"mbpd_ratio"},
            {"field":"mbc_ratio"},
            {"field":"mtwav"},
            {"field":"lateral_acc_max"},
            {"field":"long_acc_max"},
            {"field":"speed_max_level"},
            {"field":"mtpd_ratio_level"},
            {"field":"mtc_ratio_level"},
            {"field":"mbpd_ratio_level"},
            {"field":"mbc_ratio_level"},
            {"field":"mtwav_level"},
            {"field":"lateral_acc_max_level"},
            {"field":"long_acc_max_level"},
            {"field":"create_time"},
            {"field":"create_user"},
            {"field":"update_time"},
            {"field":"update_user"},
            {"field":"version"}            
      ],
      pagination:{current:1,pageSize:1},
      filter:{"id":{"Op.eq":id}}
    }
  
    let options = {headers,body:JSON.stringify(body),method:"POST",cache:'no-store'}
    const res = await fetch(`${process.env.CRV_SERVICE_URL}/data/query`,options)
    // The return value is *not* serialized
    // You can return Date, Map, Set, etc.
   
    if (!res.ok) {
      // This will activate the closest `error.js` Error Boundary
      throw new Error('Failed to fetch data')
    }
   
    return res.json()
  }