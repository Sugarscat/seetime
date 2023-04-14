import React from 'react';
import {Route, Routes} from 'react-router-dom'
import {useState} from "react";

import ElHeader from "./components/el_header";
import Login from "./pages/login";
import Home from "./pages/home";

import './App.css';

function App() {
    let [winHeight, updateWinHeight] = useState(window.innerHeight)
    const [headerH, setHeader] = useState(40)

    window.addEventListener('resize', function (e) {
        if (window.innerHeight > 512){
            updateWinHeight(window.innerHeight)
        }
    }, false);

    return (
        <div className={"App"}>
            <header className={"App-header"} style={{height:headerH}}>
                <ElHeader/>
            </header>
            <main className={"App-main"} style={{height:winHeight-headerH-0.1}}>
                <Routes>
                    <Route path="/login" element={<Login/>}></Route>
                    <Route path="/" element={<Home/>}></Route>
                </Routes>
            </main>
        </div>
    );
}

export default App;
