import request from '@/utils/request'

export function serviceList(query) {
    return request({
        url: '/service/service_list',
        method: 'get',
        params: query
    })
}

export function serviceDelete(query) {
    return request({
        url: '/service/service_delete',
        method: 'get',
        params: query
    })
}

export function serviceAddHTTP(data) {
    return request({
        url: '/service/service_add_http',
        method: 'post',
        data
    })
}

export function serviceUpdateHTTP(data) {
    return request({
        url: '/service/service_update_http',
        method: 'post',
        data
    })
}

export function serviceDetail(query) {
    return request({
        url: '/service/service_detail',
        method: 'get',
        params: query
    })
}
