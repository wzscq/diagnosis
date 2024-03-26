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