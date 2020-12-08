/* eslint-disable */
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

export function serviceAddTCP(data) {
    return request({
        url: '/service/service_add_tcp',
        method: 'post',
        data
    })
}

export function serviceAddGRPC(data) {
    return request({
        url: '/service/service_add_grpc',
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

export function serviceUpdateTCP(data) {
    return request({
        url: '/service/service_update_tcp',
        method: 'post',
        data
    })
}

export function serviceUpdateGRPC(data) {
    return request({
        url: '/service/service_update_grpc',
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

export function serviceStat(query) {
    return request({
        url: '/service/service_stat',
        method: 'get',
        params: query
    })
}
