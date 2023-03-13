package driver

//
//type traffic struct {
//	apintoDriverName string
//}
//
//func CreateTraffic(apintoDriverName string) IStrategyDriver {
//	return &traffic{apintoDriverName: apintoDriverName}
//}
//
//func (t *traffic) CheckInput(input *dto.StrategyInfoInput) error {
//	input.Uuid = strings.TrimSpace(input.Uuid)
//	if input.Uuid != "" {
//		err := common.IsMatchString(common.UUIDExp, input.Uuid)
//		if err != nil {
//			return err
//		}
//	}
//
//	input.Name = strings.TrimSpace(input.Name)
//	if input.Name == "" {
//		return errors.New("Name can't be null. ")
//	}
//	if input.Priority < 0 {
//		input.Priority = 0
//	}
//
//	if input.Config == nil {
//		return errors.New("config can't be null. ")
//	}
//
//	filterNameSet := make(map[string]struct{})
//	for _, filter := range input.Filters {
//		switch filter.Name {
//		case enum.FilterApplication, enum.FilterApi, enum.FilterPath, enum.FilterService, enum.FilterMethod:
//		default:
//			if !common.IsMatchFilterAppKey(filter.Name) {
//				return fmt.Errorf("filter.Name %s is illegal. ", filter.Name)
//			}
//		}
//
//		if len(filter.Values) == 0 {
//			return fmt.Errorf("filter.Options can't be null. filter.Name:%s ", filter.Name)
//		}
//
//		if _, has := filterNameSet[filter.Name]; has {
//			return fmt.Errorf("filterName %s is reduplicative. ", filter.Name)
//		}
//		filterNameSet[filter.Name] = struct{}{}
//	}
//
//	return nil
//}
//
//func (t *traffic) ToApinto(name, desc string, isStop bool, priority int, filters []entry.StrategyFiltersConfig, limit entry.StrategyTrafficLimitConfig) *v1.StrategyInfo {
//
//	limitingFilters := make(map[string][]string)
//
//	for _, filter := range filters {
//		limitingFilters[filter.Name] = filter.Values
//	}
//
//	return &v1.StrategyInfo{
//		Name:     name,
//		Stop:     isStop,
//		Desc:     desc,
//		Priority: priority,
//		Filters:  limitingFilters,
//		Limiting: v1.StrategyLimiting{
//			Metrics: limit.Metrics,
//			Query: v1.StrategyLimit{
//				Second: limit.Query.Second,
//				Minute: limit.Query.Minute,
//				Hour:   limit.Query.Hour,
//			},
//			Traffic: v1.StrategyLimit{
//				Second: limit.Traffic.Second,
//				Minute: limit.Traffic.Minute,
//				Hour:   limit.Traffic.Hour,
//			},
//		},
//	}
//}
