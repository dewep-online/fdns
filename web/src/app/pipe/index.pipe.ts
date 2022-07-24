import {Pipe, PipeTransform} from '@angular/core';
import {Fixed} from 'src/app/models/fixed';

@Pipe({
    name: 'index'
})
export class IndexPipe implements PipeTransform {

    transform(list: Fixed[], index: number): Fixed[] {
        return list.filter((value, i) => {
            return index === i;
        });
    }

}
