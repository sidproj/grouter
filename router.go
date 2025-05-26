package router

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)


type ctxKey string

const pathParamsKey ctxKey = "pathParams"

type handlers struct{
	get func(w http.ResponseWriter,r *http.Request)
	post func(w http.ResponseWriter,r *http.Request)
	put func(w http.ResponseWriter,r *http.Request)
	delete func(w http.ResponseWriter,r *http.Request)
}

func (h* handlers)getHandler(method string) func(http.ResponseWriter,*http.Request){
	switch method{
		case http.MethodGet: return h.get
		case http.MethodPost: return h.post
		case http.MethodPut: return h.put
		case http.MethodDelete:return h.delete
	}
	return nil
}

func GetPathParams(r *http.Request)map[string]string{
	return r.Context().Value(pathParamsKey).(map[string]string)
}

type RouterNode struct{
	handlers handlers 
	path string
	childRouters map[string]*RouterNode 
	// the about member's key for the map can be param. to know if it is dynamic route, the node's path would be empty
}

func Get(path string,handler func(w http.ResponseWriter,r *http.Request)){
	rootNode.AddChildRouters(path,"GET",handler)
}

func Post(path string,handler func(w http.ResponseWriter,r *http.Request)){
	rootNode.AddChildRouters(path,"POST",handler)
}

func Put(path string,handler func(w http.ResponseWriter,r *http.Request)){
	rootNode.AddChildRouters(path,"PUT",handler)
}

func Delete(path string,handler func(w http.ResponseWriter, r *http.Request)){
	rootNode.AddChildRouters(path,"DELETE",handler)
}

func (n* RouterNode)findHandler(path string)(*RouterNode, map[string]string){
	pathSlice := strings.Split(path, "/")[1:]
	var finalNode * RouterNode
	params := map[string]string{}
	for _,pathSeg := range pathSlice{
		var curNode * RouterNode = nil
		if(len(n.childRouters) == 0){
			return nil,nil
		}
		for k,v:= range n.childRouters{
			if k == pathSeg{
				curNode = v
				n = v
				break
			}
		}
		if curNode == nil {
			// check if pathSeg is here
			for k,v:= range n.childRouters{
				if strings.HasPrefix(k,":") {
					key := strings.TrimPrefix(k,":")
					curNode = v
					params[key] = pathSeg
					n = v
					break
				}
			}
			if curNode == nil{
				return curNode,nil
			}
		}
		finalNode = curNode
	}
	return finalNode,params;
}

func (r * RouterNode)AddChildRouters(path string,method string,handler func(w http.ResponseWriter,r *http.Request)){
	if path == "/"{
		switch method{
			case http.MethodGet: r.handlers.get = handler
			case http.MethodPost: r.handlers.post = handler
			case http.MethodPut: r.handlers.put = handler
			case http.MethodDelete: r.handlers.delete = handler
		}
		return
	}
	pathSlices := strings.Split(path,"/")[1:]
	var travelNode = r
	for _,pathSeg := range pathSlices{
		if _,ok := travelNode.childRouters[pathSeg];!ok{
			travelNode.childRouters[pathSeg] = &RouterNode{
				handlers:handlers{},
				path:"/"+pathSeg,
				childRouters: make(map[string]*RouterNode),
			}
		}
		travelNode = travelNode.childRouters[pathSeg]
	}

	switch method{
		case http.MethodGet: travelNode.handlers.get = handler
		case http.MethodPost: travelNode.handlers.post = handler
		case http.MethodPut: travelNode.handlers.put = handler
		case http.MethodDelete: travelNode.handlers.delete = handler
	}
}

var rootNode = RouterNode{
	path: "/",
	handlers: handlers{},
    childRouters: make(map[string]*RouterNode),
}

// provides mapping for all the requests
func wrapper() func (w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
		if r.URL.Path == ""{
			r.URL.Path = "/"
		}
		fmt.Printf("Request: %s , method: %s\n",r.URL.Path,r.Method)
		
		node := &rootNode
		if r.URL.Path!="/"{
			fmt.Println("Started finding node...")
			findNode,params := rootNode.findHandler(r.URL.Path)
		
			if(findNode==nil){
				http.ServeFile(w,r,"views\\404.html")
				return	
			}
			node = findNode
			if(params != nil){
				ctx := context.WithValue(r.Context(),pathParamsKey,params)
				r = r.WithContext(ctx)
			}
		}

		if(r.Method == http.MethodOptions){
			requestedMethod := r.Header.Get("Access-Control-Request-Method")
			if h := node.handlers.getHandler(requestedMethod); h!= nil{
				h(w,r)
			}
			return
		}

		if h:= node.handlers.getHandler(r.Method); h!=nil{
			h(w,r)
			return
		}

		fmt.Fprintf(w,"Invalid route")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func LoadRoutes(){
	http.HandleFunc("/",wrapper())
}