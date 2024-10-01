import React from 'react'
import { reducer } from "./reducer";
import { initialState } from './index';

const AppStateContext = React.createContext();

export const Provider = ({ children }) => {
    const [state, dispatch] = React.useReducer(reducer, initialState);

    return (
        <AppStateContext.Provider value={{ state, dispatch }}>
            {children}
        </AppStateContext.Provider>
    );
};

export const useStateValue = () => React.useContext(AppStateContext);
