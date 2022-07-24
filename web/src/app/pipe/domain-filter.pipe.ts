import {Pipe, PipeTransform} from '@angular/core';
import {AdblockDomain, AdblockURI} from 'src/app/models/adblock';

@Pipe({
    name: 'domainFilter'
})
export class DomainFilterPipe implements PipeTransform {

    transform(list: AdblockDomain[], arg: AdblockURI, filter: string): AdblockDomain[] {
        return list.filter(value => {
            if (value.tag !== arg.tag) {
                return false;
            }
            if (filter.length < 3) {
                return false;
            }
            return value.domain.indexOf(filter) !== -1;
        });
    }

}
