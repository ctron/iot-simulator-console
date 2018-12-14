import React from "react";
import {
    AggregateStatusCount,
    Card,
    CardTitle,
    CardGrid,
    Icon,
    Row,
    Col,
    AggregateStatusNotifications, AggregateStatusNotification
} from "patternfly-react";

import "./App.css"
import CardBody from "patternfly-react/dist/esm/components/Cards/CardBody";

class App extends React.Component {
    render() {
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
