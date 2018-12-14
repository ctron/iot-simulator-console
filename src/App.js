import React from "react";
import {
    AggregateStatusCount,
    Card,
    CardTitle,
    CardGrid,
    CardBody,
    Icon,
    Row,
    Col,
    AggregateStatusNotifications,
    AggregateStatusNotification,
    Masthead,
    MastheadCollapse,
    MastheadDropdown,
    MenuItem,
} from "patternfly-react";

import "./App.css"

import HorizontalNav from "patternfly-react/dist/esm/components/HorizontalNav/HorizontalNav";
import HorizontalNavMenu from "patternfly-react/dist/esm/components/HorizontalNav/HorizontalNavMenu";
import HorizontalNavMenuItem from "patternfly-react/dist/esm/components/HorizontalNav/HorizontalNavMenuItem";

class App extends React.Component {
    render() {
        return <HorizontalNav>
            <Masthead
                title="IoT simulator console"
                navToggle={false}
                thin
            />
            <HorizontalNavMenu>
                <HorizontalNavMenuItem id="overview" title="Overview"></HorizontalNavMenuItem>
            </HorizontalNavMenu>
        </HorizontalNav>
        return <Home/>
    }
}

export default App;

class Home extends React.Component {
    render() {
        return (
            <div className="cards-pf">
                <CardGrid>
                    <Row>
                        <Col>
                            <Card accented aggregated>
                                <CardTitle>
                                    <a href="#">
                                        <Icon name="shield"/>
                                        <AggregateStatusCount>7</AggregateStatusCount>
                                    </a>
                                </CardTitle>
                                <CardBody>
                                    <AggregateStatusNotifications>
                                        <AggregateStatusNotification>
                                            <Icon type="pf" name="add-circle-o"/>
                                        </AggregateStatusNotification>
                                    </AggregateStatusNotifications>
                                </CardBody>
                            </Card>
                        </Col>
                    </Row>
                </CardGrid>
            </div>
        )
    }
}
