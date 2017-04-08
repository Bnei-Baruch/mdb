import { createStore as _createStore, compose } from 'redux';
import rootReducer from './modules';

export default function createStore(data = {}) {
  const enhancers = [];

  if (process.env.NODE_ENV !== 'production' && window.__REDUX_DEVTOOLS_EXTENSION__) {
      enhancers.push(window.__REDUX_DEVTOOLS_EXTENSION__());
  }

  const store = compose(...enhancers)(_createStore)(rootReducer, data);

  return store;
}