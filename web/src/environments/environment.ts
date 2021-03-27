
import { extend } from 'lodash';
import { urls } from './urls';

export const environment = extend(urls, {
  production: false,
});
